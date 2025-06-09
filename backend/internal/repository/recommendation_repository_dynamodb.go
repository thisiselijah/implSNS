// internal/repository/recommendation_repository_dynamodb.go
package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const UserRecommendationsTableName = "UserRecommendations"

// RecommendationRepository 定義了推薦項目的操作
type RecommendationRepository interface {
	SaveRecommendations(ctx context.Context, recommendations []models.UserRecommendationItem) error
	GetRecommendations(ctx context.Context, userID string, limit int32) ([]models.UserRecommendationItem, error)
}

type dynamoDBRecommendationRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBRecommendationRepository 是 dynamoDBRecommendationRepository 的建構子
func NewDynamoDBRecommendationRepository(client *dynamodb.Client) RecommendationRepository {
	return &dynamoDBRecommendationRepository{
		client:    client,
		tableName: UserRecommendationsTableName,
	}
}

// SaveRecommendations 使用 BatchWriteItem 批量儲存推薦項目
func (r *dynamoDBRecommendationRepository) SaveRecommendations(ctx context.Context, recommendations []models.UserRecommendationItem) error {
	if len(recommendations) == 0 {
		return nil
	}

	writeRequests := make([]types.WriteRequest, len(recommendations))
	for i, item := range recommendations {
		av, err := attributevalue.MarshalMap(item)
		if err != nil {
			log.Printf("failed to marshal recommendation item for user %s: %v", item.PK, err)
			return fmt.Errorf("failed to marshal recommendation item: %w", err)
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
			log.Printf("failed to batch write recommendation items: %v", err)
			return fmt.Errorf("failed to batch write recommendation items: %w", err)
		}
	}

	log.Printf("Successfully saved %d recommendations.", len(recommendations))
	return nil
}

// GetRecommendations 獲取指定使用者的推薦列表，按分數（SK）降序排列
func (r *dynamoDBRecommendationRepository) GetRecommendations(ctx context.Context, userID string, limit int32) ([]models.UserRecommendationItem, error) {
	pkValue := "USER#" + userID

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pkValue},
		},
		ScanIndexForward: aws.Bool(false), // 根據 SK (分數) 降序排列
		Limit:            aws.Int32(limit),
	}

	result, err := r.client.Query(ctx, queryInput)
	if err != nil {
		log.Printf("DynamoDB Query failed for user recommendations %s: %v", userID, err)
		return nil, fmt.Errorf("failed to query user recommendations: %w", err)
	}

	var recommendations []models.UserRecommendationItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &recommendations); err != nil {
		log.Printf("Failed to unmarshal recommendations for user %s: %v", userID, err)
		return nil, err
	}

	return recommendations, nil
}