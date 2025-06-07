// internal/models/post_model.go (或者 feed_model.go)
package models

// FeedItem 代表在 DynamoDB 中儲存的 Feed 項目
type FeedItem struct {
	PK                      string    `dynamodbav:"PK"`
	SK                      string    `dynamodbav:"SK"`
	EntityType              string    `dynamodbav:"entity_type"`
	PostID                  string    `dynamodbav:"post_id"` // 指向原始 Post 的 ID
	AuthorID                string    `dynamodbav:"author_id"` // 原始 Post 的作者 ID
	OriginalPostCreatedAt   string    `dynamodbav:"original_post_created_at"` // 原始 Post 的創建時間 (字串格式，用於排序)
	FeedOwnerID             string    `dynamodbav:"feed_owner_id"`
}


type PostFeedDTO struct {
	PostID       string      `json:"post_id"`
	AuthorID     string      `json:"author_id"` // 原始 Post 的作者 ID
	AuthorName   string      `json:"author_name"` // 從 AuthorID 查詢得到的使用者名稱
	Content      string      `json:"content"`
	Media        []MediaItem `json:"media,omitempty"`      // MediaItem 應已在 feed_model.go 中定義
	Tags         []string    `json:"tags,omitempty"`       // stringset 在 DynamoDB, JSON 為 array of strings
	Location     *Location   `json:"location,omitempty"`   // Location 應已在 feed_model.go 中定義
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	CreatedAt    string      `json:"created_at"` // ISO 8601 String
	UpdatedAt    string      `json:"updated_at"` // ISO 8601 String
}

// 注意: models.User 結構已在 user_model.go 中定義。
// 如果 Post 結構需要嵌入 User 的部分資訊 (例如作者的用戶名)，
// 你需要在從 DynamoDB 獲取 Post 後，再額外查詢 SQL User 表來填充。
// 或者在 fan-out-on-write 時，將部分作者資訊冗餘存入 Post 實體。