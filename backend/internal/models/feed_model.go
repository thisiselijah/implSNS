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

// Post 結構體已在 user_model.go 中定義，這裡可能需要更詳細的 Post 結構
// 如果 user_model.go 中的 Post 結構不完整，或者你想分離，可以在這裡重新定義或擴展
// 例如，一個更完整的 Post 結構：
type Post struct {
	PK           string      `dynamodbav:"PK"`      // 例如 USER#{author_id}
	SK           string      `dynamodbav:"SK"`      // 例如 POST#{timestamp}#{post_id}
	GSI1PK       string      `dynamodbav:"GSI1PK"`  // 例如 POST#{post_id}
	GSI1SK       string      `dynamodbav:"GSI1SK"`  // 例如 METADATA
	EntityType   string      `dynamodbav:"entity_type"`
	PostID       string      `dynamodbav:"post_id"` // 方便直接存取
	AuthorID     string      `dynamodbav:"author_id"`
	Content      string      `dynamodbav:"content"`
	Media        []MediaItem `dynamodbav:"media,omitempty"` // omitempty 如果為空則不儲存
	Tags         []string    `dynamodbav:"tags,stringset,omitempty"` // DynamoDB String Set
	Location     *Location   `dynamodbav:"location,omitempty"`
	LikeCount    int         `dynamodbav:"like_count"`
	CommentCount int         `dynamodbav:"comment_count"`
	CreatedAt    string      `dynamodbav:"created_at"` // ISO 8601 String
	UpdatedAt    string      `dynamodbav:"updated_at"` // ISO 8601 String
}

// MediaItem 和 Location 結構也需要定義 (如果 Post 結構中使用它們)
type MediaItem struct {
	Type string `dynamodbav:"type"`
	URL  string `dynamodbav:"url"`
}
type Location struct {
	Name      string  `dynamodbav:"name"`
	Latitude  float64 `dynamodbav:"latitude"`
	Longitude float64 `dynamodbav:"longitude"`
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