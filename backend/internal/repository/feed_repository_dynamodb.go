// internal/repository/feed_repository_dynamodb.go
package repository

import (
	"backend/internal/models" // 假設有 models.UserFeedItem
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const UserFeedTableName = "UserFeed" // UserFeed 表名

type FeedRepository interface {
	GetUserFeed(ctx context.Context, userID string, limit int32, lastEvaluatedKey map[string]types.AttributeValue) (*models.PaginatedFeed, error)
	BatchAddToFeed(ctx context.Context, items []models.UserFeedItem) error // <--- 新增此方法
}

type dynamoDBFeedRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBFeedRepository(client *dynamodb.Client) FeedRepository {
	return &dynamoDBFeedRepository{
		client:    client,
		tableName: UserFeedTableName,
	}
}

func (r *dynamoDBFeedRepository) BatchAddToFeed(ctx context.Context, items []models.UserFeedItem) error {
	if len(items) == 0 {
		return nil
	}

	writeRequests := make([]types.WriteRequest, len(items))
	for i, item := range items {
		av, err := attributevalue.MarshalMap(item)
		if err != nil {
			log.Printf("failed to marshal feed item for user %s: %v", item.PK, err)
			return fmt.Errorf("failed to marshal feed item: %w", err)
		}
		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: av,
			},
		}
	}

	// DynamoDB BatchWriteItem 每次最多處理 25 個項目
	chunkSize := 25
	for i := 0; i < len(writeRequests); i += chunkSize {
		end := i + chunkSize
		if end > len(writeRequests) {
			end = len(writeRequests)
		}
		chunk := writeRequests[i:end]

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				r.tableName: chunk,
			},
		}

		_, err := r.client.BatchWriteItem(ctx, input)
		if err != nil {
			log.Printf("failed to batch write feed items: %v", err)
			return fmt.Errorf("failed to batch write feed items: %w", err)
		}
	}

	log.Printf("Successfully fanned out to %d feeds.", len(items))
	return nil
}

// GetUserFeed 從 UserFeed 表獲取 Feed
func (r *dynamoDBFeedRepository) GetUserFeed(ctx context.Context, userID string, limit int32, lastEvaluatedKey map[string]types.AttributeValue) (*models.PaginatedFeed, error) {
	pkValue := "USER#" + userID

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pkValue},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(limit),
		ExclusiveStartKey: lastEvaluatedKey,
	}

	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		log.Printf("DynamoDB Query failed for user feed %s: %v", userID, err)
		return nil, fmt.Errorf("failed to query user feed: %w", err)
	}

	// --- 修改開始：手動解析以處理不一致的 TTLTimestamp 類型 ---
	var feedItems []models.UserFeedItem
	// 遍歷從 DynamoDB 返回的每一個項目
	for _, itemMap := range result.Items {
		var feedItem models.UserFeedItem
		// 1. 先將整個項目反序列化到一個通用的 map[string]interface{} 中
		var rawMap map[string]interface{}
		if err := attributevalue.UnmarshalMap(itemMap, &rawMap); err != nil {
			log.Printf("Failed to unmarshal raw feed item map: %v. Skipping item.", err)
			continue // 跳過此筆損壞的資料
		}

		// 2. 手動將 map 中的值賦給我們的 struct 欄位
		if pk, ok := rawMap["PK"].(string); ok {
			feedItem.PK = pk
		}
		if sk, ok := rawMap["SK"].(string); ok {
			feedItem.SK = sk
		}
		if postID, ok := rawMap["PostID"].(string); ok {
			feedItem.PostID = postID
		}
		if authorID, ok := rawMap["AuthorID"].(string); ok {
			feedItem.AuthorID = authorID
		}

		// 3. 彈性處理 TTLTimestamp 欄位
		if ttlVal, ok := rawMap["TTLTimestamp"]; ok {
			switch v := ttlVal.(type) {
			case float64: // DynamoDB Number 在 Go 中通常被解析為 float64
				feedItem.TTLTimestamp = int64(v)
			case string:
				// 如果是字串，我們嘗試將其轉換為 int64
				parsedTTL, parseErr := strconv.ParseInt(v, 10, 64)
				if parseErr == nil {
					feedItem.TTLTimestamp = parsedTTL
				} else {
					log.Printf("Could not parse TTLTimestamp string '%s' to int64. TTL will be 0 for this item.", v)
				}
			}
		}
		feedItems = append(feedItems, feedItem)
	}
	// --- 修改結束 ---

	paginatedFeed := &models.PaginatedFeed{
		Items:            feedItems,
		LastEvaluatedKey: result.LastEvaluatedKey,
	}

	return paginatedFeed, nil
}
