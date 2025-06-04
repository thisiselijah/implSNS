// internal/repository/post_repository_dynamodb.go
package repository

import (
	"context"
	"fmt"
	"log"
	// "sort" // 用於排序
	"strings"
	"sync"

	"backend/internal/models" // 假設您有 models.FeedItem 和 models.Post 結構

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const FeedTableName = "Posts" // 假設您的表名

// DynamoDBPostRepository 結構
type DynamoDBPostRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBPostRepository 建構子
func NewDynamoDBPostRepository(client *dynamodb.Client) *DynamoDBPostRepository {
	return &DynamoDBPostRepository{
		client:    client,
		tableName: FeedTableName, // 或者從配置讀取
	}
}

// GetFeedItemsByUserID 從 DynamoDB 獲取指定用戶的 Feed Item 列表
// PK = USER#{userID}, SK starts_with FEEDITEM#
// 返回的 feedItems 應該按 SK (時間戳) 排序
func (r *DynamoDBPostRepository) GetFeedItemsByUserID(ctx context.Context, userPK string) ([]models.FeedItem, error) {
	// userPK 應該是 "USER#userID_actual_value"
	log.Printf("Fetching feed items for PK: %s", userPK)

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pkval AND begins_with(SK, :skprefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval":    &types.AttributeValueMemberS{Value: userPK},
			":skprefix": &types.AttributeValueMemberS{Value: "FEEDITEM#"},
		},
		ScanIndexForward: aws.Bool(false), // false 表示按 SK 降序 (最新的在前)
		// Limit: aws.Int32(20), // 可以添加分頁限制
	}

	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		log.Printf("DynamoDB Query failed for feed items PK %s: %v", userPK, err)
		return nil, fmt.Errorf("failed to query feed items: %w", err)
	}

	var feedItems []models.FeedItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &feedItems); err != nil {
		log.Printf("Failed to unmarshal feed items for PK %s: %v", userPK, err)
		return nil, fmt.Errorf("failed to unmarshal feed items: %w", err)
	}
	log.Printf("Successfully fetched %d feed items for PK: %s", len(feedItems), userPK)


	// 由於 DynamoDB 的 begins_with 和 ScanIndexForward 已經排序，這裡通常不需要額外排序
	// 如果需要基於 FeedItem 結構中的特定時間戳欄位（例如 OriginalPostCreatedAt）再次確認排序，可以在這裡做
	// sort.SliceStable(feedItems, func(i, j int) bool {
	// 	return feedItems[i].OriginalPostCreatedAt > feedItems[j].OriginalPostCreatedAt // 假設是字串且可比較，或轉換為 time.Time
	// })


	return feedItems, nil
}

// GetPostsByIDs 從 DynamoDB 批量獲取完整的 Post 內容
// 這通常使用 BatchGetItem，或者如果 postID 都是 GSI1PK，則可以多次查詢 GSI
// 這裡我們示範使用 BatchGetItem
func (r *DynamoDBPostRepository) GetPostsByIDs(ctx context.Context, postIDs []string) ([]models.Post, error) {
	if len(postIDs) == 0 {
		return []models.Post{}, nil
	}

	log.Printf("Fetching posts by IDs: %v", postIDs)

	keys := make([]map[string]types.AttributeValue, len(postIDs))
	for i, id := range postIDs {
		// 假設 Post 的主資料表 PK 是 USER#{author_id}, SK 是 POST#{timestamp}#{post_id}
		// 而我們通常使用 GSI1 (GSI1PK: POST#{post_id}, GSI1SK: METADATA) 來直接獲取貼文
		// 如果 BatchGetItem 用於主表，你需要知道每個 postID 對應的完整主鍵 (PK 和 SK)
		// 更常見的做法是，如果 Post 實體是透過 GSI 查詢的，
		// BatchGetItem 也可以用在 GSI 上，但這裡的 key 結構會不同。

		// 假設我們用 BatchGetItem 查詢 GSI1 (PostLookup)
		// GSI1PK: POST#{post_id}, GSI1SK: METADATA
		// BatchGetItem 需要的是基礎表的主鍵。所以這意味著我們需要一種方法來從 postID 映射回基礎表的主鍵
		// 或者，如果 GetPostByID 內部是查詢 GSI，我們可以多次調用 GetPostByID。
		// 但為了效率，BatchGetItem 是首選。

		// 簡化：這裡我們假設 GSI1 投影了所有需要的屬性，並且我們能構造出基礎表的主鍵
		// 或者，更簡單的方式是，如果 Post 實體自身就是以 POST#{post_id} 作為 PK，
		// SK 是 METADATA (或者 POST#{post_id})，那麼 BatchGetItem 就很直接。

		// 在我們之前的設計中，Post 的 PK 是 USER#{author_id}, SK 是 POST#{timestamp}#{post_id}
		// GSI1PK 是 POST#{post_id}
		// BatchGetItem 操作作用於基礎表。所以我們不能直接用 BatchGetItem 和 ["POST#id1", "POST#id2"]
		// 除非我們知道每個 post_id 對應的完整主鍵。

		// 因此，這裡有幾種策略：
		// 1. Fan-out-on-write 時，FeedItem 也儲存 Post 的 PK 和 SK (冗餘，但讀取快)。
		// 2. 多次並行執行 GetItem on GSI1 (每次查一個 post_id)。
		// 3. 如果 DynamoDB 表設計允許，並且 Post 有一個簡單的 PK (如 post_id)，則可以直接用 BatchGetItem。

		// 這裡採用策略 2 的簡化版：多次查詢 GSI1
		// 注意：在生產環境中，應使用並行查詢（例如 goroutines + channels）來提高效率
		// 這裡為了簡單，串行查詢，但這不是最佳實踐。
		keys[i] = map[string]types.AttributeValue{
			// 這是 GSI1 (PostLookup) 的鍵結構
			"GSI1PK": &types.AttributeValueMemberS{Value: "POST#" + id},
			"GSI1SK": &types.AttributeValueMemberS{Value: "METADATA"}, // 或者 PostID，取決於 GSI1SK 的設計
		}
	}


	// 策略 2 的實際執行：多次並行 Query GSI1
	// 由於 BatchGetItem 直接用於 GSI 是不支援的，我們需要多次 Query GSI 或 GetItem GSI
	var posts []models.Post
	var mu sync.Mutex // Mutex to protect posts slice
    var wg sync.WaitGroup

	for _, id := range postIDs {
		wg.Add(1)
		go func(postID string) {
			defer wg.Done()
			// 這裡模擬 GetPostByID 內部查詢 GSI1 的邏輯
			queryInput := &dynamodb.QueryInput{
				TableName:              aws.String(r.tableName),
				IndexName:              aws.String("GSI1"), // 使用 GSI1 (PostLookup)
				KeyConditionExpression: aws.String("GSI1PK = :gsi1pkval AND GSI1SK = :gsi1skval"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":gsi1pkval":    &types.AttributeValueMemberS{Value: "POST#" + postID},
					":gsi1skval": &types.AttributeValueMemberS{Value: "METADATA"},
				},
			}
			result, err := r.client.Query(ctx, queryInput)
			if err != nil {
				log.Printf("Error querying GSI1 for postID %s: %v", postID, err)
				return
			}
			if len(result.Items) > 0 {
				var post models.Post
				if err := attributevalue.UnmarshalMap(result.Items[0], &post); err == nil {
					// 確保 PostID 被正確填充
					if post.PostID == "" && strings.HasPrefix(post.GSI1PK, "POST#") { // 假設 GSI1PK 包含 PostID
						post.PostID = strings.TrimPrefix(post.GSI1PK, "POST#")
					}
					mu.Lock()
					posts = append(posts, post)
					mu.Unlock()
				} else {
					log.Printf("Error unmarshalling post for postID %s: %v", postID, err)
				}
			} else {
                log.Printf("Post not found via GSI1 for postID %s", postID)
            }
		}(id)
	}
    wg.Wait()


	log.Printf("Successfully fetched %d posts by IDs", len(posts))
	return posts, nil
}


// --- 其他 Post Repository 方法 (CreatePost, GetPostByID 等) ---
// func (r *DynamoDBPostRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) { ... }
// func (r *DynamoDBPostRepository) CreatePost(ctx context.Context, post *models.Post) error { ... }
// ...