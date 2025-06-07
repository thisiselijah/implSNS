// backend/internal/handler/post_handler.go
package handler

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"
	// "fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// PostHandler 結構
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler 建構子
func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// getAuthenticatedUserID 是一個輔助函式，從 Gin context 中獲取並驗證使用者 ID
func getAuthenticatedUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token context"})
		return "", false
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID in context is of invalid type or empty"})
		return "", false
	}
	return userIDStr, true
}

// CreatePost 處理新增貼文請求
func (h *PostHandler) CreatePost(c *gin.Context) {
	var payload models.CreatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 從 context 獲取已驗證的使用者 ID，並覆寫 payload 中的 AuthorID
	authorID, ok := getAuthenticatedUserID(c)
	log.Printf("Authenticated user ID: %s", authorID) // 日誌輸出，便於調試
	if !ok {
		return // 錯誤已由輔助函式發送
	}
	payload.AuthorID = authorID

	post, err := h.postService.CreatePost(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPostsByUserID 處理獲取作者貼文的請求
func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	posts, err := h.postService.GetPostsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// UpdatePost 處理編輯貼文的請求
// 注意：此處的 service.UpdatePost 方法在原始設計中未包含權限驗證。
// 一個更安全的設計是讓 service 層接收 authenticatedUserID 並在更新前進行比對。
func (h *PostHandler) UpdatePost(c *gin.Context) {
	var payload models.UpdatePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 在理想情況下，service 層應該執行權限檢查。
	// _, ok := getAuthenticatedUserID(c)
	// if !ok {
	// 	return
	// }
	// 建議修改 `UpdatePost` 服務，使其能接收並驗證操作者 ID

	updatedPost, err := h.postService.UpdatePost(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    updatedPost,
	})
}

// DeletePost 處理刪除貼文的請求
func (h *PostHandler) DeletePost(c *gin.Context) {
	var payload models.DeletePostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// 從 context 獲取已驗證的使用者 ID，並覆寫 payload 中的 AuthorID，交由 service 層進行權限驗證
	authorID, ok := getAuthenticatedUserID(c)
	if !ok {
		return // 錯誤已由輔助函式發送
	}
	payload.AuthorID = authorID

	err := h.postService.DeletePost(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}


// GetFeedPosts 處理獲取使用者 Feed 的請求
// ... (此函式主要為讀取操作，保持不變) ...
func GetFeedPosts(dynamoDBClient *dynamodb.Client, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID") // 從 URL 路徑中獲取 userID

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
			return
		}
		postRepo := repository.NewDynamoDBPostRepository(dynamoDBClient)
		feedItems, err := postRepo.GetFeedItemsByUserID(c.Request.Context(), "USER#"+userID)
		if err != nil {
			log.Printf("Error fetching feed items for userID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed items"})
			return
		}
		if len(feedItems) == 0 {
			c.JSON(http.StatusOK, []models.PostFeedDTO{})
			return
		}
		var postIDs []string
		for _, item := range feedItems {
			if item.PostID != "" {
				postIDs = append(postIDs, item.PostID)
			}
		}
		if len(postIDs) == 0 {
			log.Printf("No valid postIDs found in feed items for userID %s", userID)
			c.JSON(http.StatusOK, []models.PostFeedDTO{})
			return
		}
		postsFromDB, err := postRepo.GetPostsByIDs(c.Request.Context(), postIDs)
		if err != nil {
			log.Printf("Error fetching full posts for userID %s, postIDs %v: %v", userID, postIDs, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch full posts for feed"})
			return
		}
		var feedDTOs []models.PostFeedDTO
		postMap := make(map[string]models.Post)
		for _, p := range postsFromDB {
			postMap[p.PostID] = p
		}
		for _, feedItem := range feedItems {
			if post, ok := postMap[feedItem.PostID]; ok {
				var authorName string
				authorIDUint, convErr := strconv.ParseUint(post.AuthorID, 10, 64)
				if convErr != nil {
					log.Printf("Error converting AuthorID string '%s' to uint for post '%s': %v. Using original AuthorID as name.", post.AuthorID, post.PostID, convErr)
					authorName = "AuthorID: " + post.AuthorID
				} else {
					user, userErr := userRepo.GetUserByID(uint(authorIDUint))
					if userErr != nil {
						log.Printf("Error fetching user (ID: %d) for post '%s': %v. Using original AuthorID as name.", authorIDUint, post.PostID, userErr)
						authorName = "User: " + post.AuthorID
					} else {
						authorName = user.Username
					}
				}
				dto := models.PostFeedDTO{
					PostID:       post.PostID,
					AuthorID:     post.AuthorID,
					AuthorName:   authorName,
					Content:      post.Content,
					Media:        post.Media,
					Tags:         post.Tags,
					Location:     post.Location,
					LikeCount:    post.LikeCount,
					CommentCount: post.CommentCount,
					CreatedAt:    post.CreatedAt,
					UpdatedAt:    post.UpdatedAt,
				}
				feedDTOs = append(feedDTOs, dto)
			} else {
				log.Printf("Post with ID %s found in feedItems but not in fetched postsMap for userID %s", feedItem.PostID, userID)
			}
		}
		c.JSON(http.StatusOK, feedDTOs)
	}
}


func (h *PostHandler) LikePost(c *gin.Context) {
	postID := c.Param("postID")

	// 從 context 獲取已驗證的使用者 ID
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return // 錯誤已由輔助函式發送
	}

	if err := h.postService.LikePost(c.Request.Context(), postID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

// UnlikePost 處理取消按讚請求
func (h *PostHandler) UnlikePost(c *gin.Context) {
	postID := c.Param("postID")

	// 從 context 獲取已驗證的使用者 ID
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return // 錯誤已由輔助函式發送
	}

	if err := h.postService.UnlikePost(c.Request.Context(), postID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// CreateComment 處理新增評論請求
func (h *PostHandler) CreateComment(c *gin.Context) {
	var payload models.CreateCommentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	payload.PostID = c.Param("postID")

	// 從 context 獲取已驗證的使用者 ID，並覆寫 payload 中的 AuthorID
	authorID, ok := getAuthenticatedUserID(c)
	if !ok {
		return // 錯誤已由輔助函式發送
	}
	payload.AuthorID = authorID

	comment, err := h.postService.CreateComment(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// DeleteComment 處理刪除評論請求
func (h *PostHandler) DeleteComment(c *gin.Context) {
	postID := c.Param("postID")
	commentSK := c.Param("commentSK")

	// 從 context 獲取已驗證的使用者 ID
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		return // 錯誤已由輔助函式發送
	}

	if err := h.postService.DeleteComment(c.Request.Context(), postID, commentSK, userID); err != nil {
		// 更精確地處理權限錯誤
		if err.Error() == "user not authorized to delete this comment" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}