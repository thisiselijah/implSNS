// backend/internal/repository/post_repository_dynamodb.go
// internal/repository/post_repository_dynamodb.go
package repository

import (
	"backend/internal/models" // 假設您有 models.FeedItem 和 models.Post 結構
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type PostRepository interface {
	GetFeedItemsByUserID(ctx context.Context, userPK string) ([]models.FeedItem, error)
	GetPostsByIDs(ctx context.Context, postIDs []string) ([]models.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) error
	UpdatePost(ctx context.Context, post *models.Post) error
	DeletePost(ctx context.Context, authorID, postID, createdAt string) error
	GetPostByID(ctx context.Context, postID string) (*models.Post, error)
	GetRecentPosts(ctx context.Context, lookbackDays int) ([]models.Post, error)

	// --- FIX: Signatures changed to accept *models.Post ---
	AddLike(ctx context.Context, post *models.Post, userID string) error
	RemoveLike(ctx context.Context, post *models.Post, userID string) error
	CreateComment(ctx context.Context, post *models.Post, comment *models.Comment) error
	DeleteComment(ctx context.Context, post *models.Post, commentSK string) error
	GetCommentBySK(ctx context.Context, postID, commentSK string) (*models.Comment, error)
	CheckIfPostsLikedBy(ctx context.Context, postIDs []string, userID string) (map[string]bool, error) // <--- 新增此方法

}

const FeedTableName = "Posts" // 假設您的表名

// DynamoDBPostRepository 結構
type DynamoDBPostRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBPostRepository(client *dynamodb.Client) PostRepository {
	return &DynamoDBPostRepository{
		client:    client,
		tableName: FeedTableName, // 或者從配置讀取
	}
}

func (r *DynamoDBPostRepository) GetRecentPosts(ctx context.Context, lookbackDays int) ([]models.Post, error) {
	cutOffDate := time.Now().UTC().AddDate(0, 0, -lookbackDays).Format(time.RFC3339)

	// 使用 Scan 操作篩選近期貼文。這在大型表上效率低下。
	// 生產環境應建立 GSI (例如 PK: EntityType, SK: CreatedAt) 來高效查詢。
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("entity_type = :type AND created_at >= :date"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":type": &types.AttributeValueMemberS{Value: "POST"},
			":date": &types.AttributeValueMemberS{Value: cutOffDate},
		},
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		log.Printf("Failed to scan for recent posts: %v", err)
		return nil, err
	}

	var posts []models.Post
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &posts); err != nil {
		log.Printf("Failed to unmarshal recent posts: %v", err)
		return nil, err
	}
	return posts, nil
}

func (r *DynamoDBPostRepository) CheckIfPostsLikedBy(ctx context.Context, postIDs []string, userID string) (map[string]bool, error) {
	if len(postIDs) == 0 {
		return make(map[string]bool), nil
	}

	// 初始化結果 map，預設所有貼文都未被按讚
	likedStatus := make(map[string]bool, len(postIDs))
	for _, id := range postIDs {
		likedStatus[id] = false
	}

	// 根據按讚的資料模型 (PK: POST#{post_id}, SK: USER#{user_id}) 建立要查詢的索引鍵
	keys := make([]map[string]types.AttributeValue, len(postIDs))
	for i, postID := range postIDs {
		keys[i] = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "POST#" + postID},
			"SK": &types.AttributeValueMemberS{Value: "USER#" + userID},
		}
	}

	// BatchGetItem 每次最多查詢 100 個項目，如果超過則需分批
	chunkSize := 100
	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}
		chunk := keys[i:end]

		input := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				r.tableName: {
					Keys: chunk,
					// 我們只關心項目是否存在，不需回傳屬性，節省讀取成本
					ProjectionExpression: aws.String("PK"),
				},
			},
		}

		result, err := r.client.BatchGetItem(ctx, input)
		if err != nil {
			log.Printf("BatchGetItem failed for checking likes: %v", err)
			return nil, err
		}

		// 如果 BatchGetItem 找到了對應的項目，表示使用者按過讚
		if responses, ok := result.Responses[r.tableName]; ok {
			for _, itemMap := range responses {
				var like struct {
					PK string `dynamodbav:"PK"`
				}
				if err := attributevalue.UnmarshalMap(itemMap, &like); err == nil {
					// 從 PK "POST#..." 中解析出 postID
					postID := strings.TrimPrefix(like.PK, "POST#")
					likedStatus[postID] = true
				}
			}
		}
	}

	return likedStatus, nil
}

func (r *DynamoDBPostRepository) AddLike(ctx context.Context, post *models.Post, userID string) error {
	like := models.Like{
		PK:         "POST#" + post.PostID,
		SK:         "USER#" + userID,
		EntityType: "LIKED_POST",
		PostID:     post.PostID,
		UserID:     userID,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	}
	likeItem, err := attributevalue.MarshalMap(like)
	if err != nil {
		return fmt.Errorf("failed to marshal like item: %w", err)
	}

	// The post's key is now taken directly from the passed-in object.
	postKey, err := attributevalue.MarshalMap(map[string]string{"PK": post.PK, "SK": post.SK})
	if err != nil {
		return fmt.Errorf("failed to marshal post key for like: %w", err)
	}

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:           aws.String(r.tableName),
					Item:                likeItem,
					ConditionExpression: aws.String("attribute_not_exists(PK)"), // Prevents duplicate likes
				},
			},
			{
				Update: &types.Update{
					TableName:        aws.String(r.tableName),
					Key:              postKey,
					UpdateExpression: aws.String("ADD like_count :inc"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":inc": &types.AttributeValueMemberN{Value: "1"},
					},
				},
			},
		},
	})

	if err != nil {
		if _, ok := err.(*types.TransactionCanceledException); ok {
			return errors.New("transaction failed, possibly already liked")
		}
		log.Printf("Error in AddLike transaction: %v", err)
		return err
	}
	return nil
}

// --- FIX: RemoveLike now uses the passed-in post's PK and SK ---
func (r *DynamoDBPostRepository) RemoveLike(ctx context.Context, post *models.Post, userID string) error {
	likePK := "POST#" + post.PostID
	likeSK := "USER#" + userID

	postKey, err := attributevalue.MarshalMap(map[string]string{"PK": post.PK, "SK": post.SK})
	if err != nil {
		return fmt.Errorf("failed to marshal post key for unlike: %w", err)
	}

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					TableName: aws.String(r.tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: likePK},
						"SK": &types.AttributeValueMemberS{Value: likeSK},
					},
					ConditionExpression: aws.String("attribute_exists(PK)"), // Ensure the like exists
				},
			},
			{
				Update: &types.Update{
					TableName:           aws.String(r.tableName),
					Key:                 postKey,
					UpdateExpression:    aws.String("ADD like_count :dec"),
					ConditionExpression: aws.String("like_count > :zero"), // Prevent negative counts
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":dec":  &types.AttributeValueMemberN{Value: "-1"},
						":zero": &types.AttributeValueMemberN{Value: "0"},
					},
				},
			},
		},
	})

	if err != nil {
		if _, ok := err.(*types.TransactionCanceledException); ok {
			return errors.New("transaction failed, possibly not liked yet or count is zero")
		}
		log.Printf("Error in RemoveLike transaction: %v", err)
		return err
	}
	return nil
}

// --- FIX: CreateComment now uses the passed-in post's PK and SK ---
func (r *DynamoDBPostRepository) CreateComment(ctx context.Context, post *models.Post, comment *models.Comment) error {
	now := time.Now().UTC()
	commentID := uuid.New().String()
	timestamp := now.Format(time.RFC3339Nano)

	comment.PK = "POST#" + post.PostID
	comment.SK = fmt.Sprintf("COMMENT#%s#%s", timestamp, commentID)
	comment.CommentID = commentID
	comment.EntityType = "COMMENT"
	comment.CreatedAt = timestamp

	commentItem, err := attributevalue.MarshalMap(comment)
	if err != nil {
		return err
	}

	postKey, err := attributevalue.MarshalMap(map[string]string{"PK": post.PK, "SK": post.SK})
	if err != nil {
		return fmt.Errorf("failed to marshal post key for comment: %w", err)
	}

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String(r.tableName),
					Item:      commentItem,
				},
			},
			{
				Update: &types.Update{
					TableName:        aws.String(r.tableName),
					Key:              postKey,
					UpdateExpression: aws.String("ADD comment_count :inc"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":inc": &types.AttributeValueMemberN{Value: "1"},
					},
				},
			},
		},
	})

	if err != nil {
		log.Printf("Error in CreateComment transaction: %v", err)
		return err
	}
	return nil
}

// --- FIX: DeleteComment now uses the passed-in post's PK and SK ---
func (r *DynamoDBPostRepository) DeleteComment(ctx context.Context, post *models.Post, commentSK string) error {
	commentPK := "POST#" + post.PostID

	postKey, err := attributevalue.MarshalMap(map[string]string{"PK": post.PK, "SK": post.SK})
	if err != nil {
		return fmt.Errorf("failed to marshal post key for delete comment: %w", err)
	}

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Delete: &types.Delete{
					TableName: aws.String(r.tableName),
					Key: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{Value: commentPK},
						"SK": &types.AttributeValueMemberS{Value: commentSK},
					},
				},
			},
			{
				Update: &types.Update{
					TableName:           aws.String(r.tableName),
					Key:                 postKey,
					UpdateExpression:    aws.String("ADD comment_count :dec"),
					ConditionExpression: aws.String("comment_count > :zero"),
					ExpressionAttributeValues: map[string]types.AttributeValue{
						":dec":  &types.AttributeValueMemberN{Value: "-1"},
						":zero": &types.AttributeValueMemberN{Value: "0"},
					},
				},
			},
		},
	})

	if err != nil {
		log.Printf("Error in DeleteComment transaction: %v", err)
		return err
	}
	return nil
}

// GetCommentBySK gets a comment by its full primary key (PK and SK)
func (r *DynamoDBPostRepository) GetCommentBySK(ctx context.Context, postID, commentSK string) (*models.Comment, error) {
	pk := "POST#" + postID
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: commentSK},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, errors.New("comment not found")
	}
	var comment models.Comment
	if err := attributevalue.UnmarshalMap(result.Item, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// CreatePost 將新貼文儲存到 DynamoDB
func (r *DynamoDBPostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	// 補全必要欄位
	now := time.Now().UTC()
	postID := uuid.New().String()
	timestamp := now.Format(time.RFC3339Nano)

	post.PostID = postID
	post.CreatedAt = timestamp
	post.UpdatedAt = timestamp
	post.PK = "USER#" + post.AuthorID
	post.SK = "POST#" + timestamp + "#" + postID
	post.GSI1PK = "POST#" + postID
	post.GSI1SK = "METADATA"
	post.EntityType = "POST"

	item, err := attributevalue.MarshalMap(post)
	if err != nil {
		log.Printf("Error marshalling post for CreatePost: %v", err)
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Error putting item to DynamoDB for CreatePost: %v", err)
		return err
	}
	return nil
}

// GetPostsByAuthorID 透過 PK 查詢作者的所有貼文
func (r *DynamoDBPostRepository) GetPostsByUserID(ctx context.Context, userID string) ([]models.Post, error) {
	pk := "USER#" + userID
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pkval AND begins_with(SK, :skprefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkval":    &types.AttributeValueMemberS{Value: pk},
			":skprefix": &types.AttributeValueMemberS{Value: "POST#"},
		},
		ScanIndexForward: aws.Bool(false), // 最新貼文在前
	}

	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		log.Printf("DynamoDB Query failed for GetPostsByAuthorID PK %s: %v", pk, err)
		return nil, err
	}

	var posts []models.Post
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &posts); err != nil {
		log.Printf("Failed to unmarshal posts for GetPostsByAuthorID PK %s: %v", pk, err)
		return nil, err
	}
	return posts, nil
}

// UpdatePost 更新貼文內容
func (r *DynamoDBPostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	// 為了更新，我們需要知道完整的 Key (PK, SK)
	// Service 層應先獲取 post，然後傳遞過來
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": post.PK,
		"SK": post.SK,
	})
	if err != nil {
		return err
	}

	updateExpression := "SET content = :c, updatedAt = :u"
	expressionAttributeValues, err := attributevalue.MarshalMap(map[string]interface{}{
		":c": post.Content,
		":u": time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	if err != nil {
		log.Printf("Error updating post in DynamoDB: %v", err)
		return err
	}

	return nil
}

// DeletePost 刪除貼文
func (r *DynamoDBPostRepository) DeletePost(ctx context.Context, authorID, postID, createdAt string) error {
	// 為了刪除，我們需要重建 SK
	// 注意：這種方法要求 createdAt 的格式必須與儲存時完全一致
	pk := "USER#" + authorID
	// SK 的重建依賴於 CreatedAt 和 PostID
	// 這是一個潛在的脆弱點，如果 SK 的生成邏輯改變，這裡也要改
	sk := fmt.Sprintf("POST#%s#%s", createdAt, postID)

	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": pk,
		"SK": sk,
	})
	if err != nil {
		log.Printf("Error marshalling key for DeletePost: %v", err)
		return err
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	})

	if err != nil {
		log.Printf("Error deleting post from DynamoDB: %v", err)
		return err
	}
	return nil
}

// GetPostByID 透過 GSI 查詢單一貼文
func (r *DynamoDBPostRepository) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :gsi1pkval AND GSI1SK = :gsi1skval"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1pkval": &types.AttributeValueMemberS{Value: "POST#" + postID},
			":gsi1skval": &types.AttributeValueMemberS{Value: "METADATA"},
		},
	}
	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		log.Printf("Error querying GSI1 for GetPostByID %s: %v", postID, err)
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("post with ID %s not found", postID)
	}

	var post models.Post
	if err := attributevalue.UnmarshalMap(result.Items[0], &post); err != nil {
		log.Printf("Error unmarshalling post for GetPostByID %s: %v", postID, err)
		return nil, err
	}
	return &post, nil
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
			post, err := r.GetPostByID(ctx, postID)
			if err != nil {
				log.Printf("Warning: Could not fetch post by ID %s (it may have been deleted): %v", postID, err)
				return
			}
			mu.Lock()
			posts = append(posts, *post)
			mu.Unlock()

		}(id)
	}
	wg.Wait()

	log.Printf("Successfully fetched %d posts by IDs", len(posts))
	return posts, nil
}
