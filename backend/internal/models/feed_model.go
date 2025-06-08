// internal/models/feed_model.go
package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

// FeedItem 代表在 DynamoDB 中儲存的 Feed 項目
type FeedItem struct {
	PK                    string `dynamodbav:"PK"`
	SK                    string `dynamodbav:"SK"`
	EntityType            string `dynamodbav:"entity_type"`
	PostID                string `dynamodbav:"post_id"`                  // 指向原始 Post 的 ID
	AuthorID              string `dynamodbav:"author_id"`                // 原始 Post 的作者 ID
	OriginalPostCreatedAt string `dynamodbav:"original_post_created_at"` // 原始 Post 的創建時間 (字串格式，用於排序)
	FeedOwnerID           string `dynamodbav:"feed_owner_id"`
}

type PostFeedDTO struct {
	PostID       string      `json:"post_id"`
	AuthorID     string      `json:"author_id"`   // 原始 Post 的作者 ID
	AuthorName   string      `json:"author_name"` // 從 AuthorID 查詢得到的使用者名稱
	Content      string      `json:"content"`
	Media        []MediaItem `json:"media,omitempty"`    // MediaItem 應已在 feed_model.go 中定義
	Tags         []string    `json:"tags,omitempty"`     // stringset 在 DynamoDB, JSON 為 array of strings
	Location     *Location   `json:"location,omitempty"` // Location 應已在 feed_model.go 中定義
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	CreatedAt    string      `json:"created_at"` // ISO 8601 String
	UpdatedAt    string      `json:"updated_at"` // ISO 8601 String
	IsLiked      bool        `json:"isLiked"`
}

// UserFeedItem 代表 UserFeed 表中的一個項目
type UserFeedItem struct {
	PK           string `dynamodbav:"PK"`           // 分割區索引鍵, e.g., "USER#2"
	SK           string `dynamodbav:"SK"`           // 排序鍵, e.g., 貼文的創建時間
	PostID       string `dynamodbav:"PostID"`
	AuthorID     string `dynamodbav:"AuthorID"`
	TTLTimestamp int64  `dynamodbav:"TTLTimestamp"` // 用於 TTL 的時間戳記
}


// PaginatedFeed 用於包含分頁結果和下一次查詢的鍵
type PaginatedFeed struct {
	Items            []UserFeedItem
	LastEvaluatedKey map[string]types.AttributeValue
}

type UserRecommendationItem struct {
	PK               string `dynamodbav:"PK"`               // e.g., USER#10
	SK               string `dynamodbav:"SK"`               // e.g., 推薦分數，用於排序
	GSI1PK           string `dynamodbav:"GSI1PK"`           // e.g., "trending-v1.0"
	GSI1SK           string `dynamodbav:"GSI1SK"`           // e.g., USER#10
	PostID           string `dynamodbav:"PostID"`           // 推薦的貼文 ID
	AlgorithmVersion string `dynamodbav:"AlgorithmVersion"` // 演算法版本
	GeneratedAt      string `dynamodbav:"GeneratedAt"`      // 推薦產生的時間
}
