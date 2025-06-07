package models

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

type CreatePostPayload struct {
	AuthorID string      `json:"author_id" binding:"required"`
	Content  string      `json:"content" binding:"required"`
	Media    []MediaItem `json:"media,omitempty"`
	Tags     []string    `json:"tags,omitempty"`
	Location *Location   `json:"location,omitempty"`
}

// UpdatePostPayload 定義了編輯貼文請求的 JSON 結構
type UpdatePostPayload struct {
	PostID  string `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	// 其他允許更新的欄位...
}

// DeletePostPayload 定義了刪除貼文請求的 JSON 結構
type DeletePostPayload struct {
	PostID   string `json:"post_id" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
	// 為了刪除 DynamoDB 項目，我們需要完整的 Primary Key (PK, SK)。
	// SK 包含時間戳，前端可能沒有。
	// 這裡的設計是讓 Service 層根據 PostID 找到 Post，再刪除。
}


// --- 新增 Like 和 Comment 相關的模型 ---

// Like 記錄了誰對哪篇貼文按讚
type Like struct {
	PK         string `dynamodbav:"PK"`      // POST#{post_id}
	SK         string `dynamodbav:"SK"`      // USER#{user_id}
	EntityType string `dynamodbav:"entity_type"`
	PostID     string `dynamodbav:"post_id"`
	UserID     string `dynamodbav:"user_id"`
	CreatedAt  string `dynamodbav:"created_at"`
}

// Comment 包含了評論的詳細資訊
type Comment struct {
	PK           string `dynamodbav:"PK"`         // POST#{post_id}
	SK           string `dynamodbav:"SK"`         // COMMENT#{timestamp}#{comment_id}
	EntityType   string `dynamodbav:"entity_type"`
	CommentID    string `dynamodbav:"comment_id"`
	PostID       string `dynamodbav:"post_id"`
	AuthorID     string `dynamodbav:"author_id"`
	AuthorName   string `dynamodbav:"author_name"` // 冗餘儲存，方便查詢
	Content      string `dynamodbav:"content"`
	CreatedAt    string `dynamodbav:"created_at"`
}

// CreateCommentPayload 定義了新增評論請求的 JSON 結構
type CreateCommentPayload struct {
	PostID   string `json:"post_id" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}