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
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	viewerID, _ := getAuthenticatedUserID(c)
	posts, err := h.postService.GetPostsByUserID(c.Request.Context(), userID, viewerID)
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


func GetFeedPosts(dynamoDBClient *dynamodb.Client, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		viewerID, ok := getAuthenticatedUserID(c)
		if !ok {
			// getAuthenticatedUserID 內部已處理錯誤回應
			return
		}
		userID := c.Param("userID")
		
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
			return
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		exclusiveStartKey, _ := c.GetQuery("next_key")

		var lastEvaluatedKey map[string]types.AttributeValue
		if exclusiveStartKey != "" {
			keyJSON, err := base64.StdEncoding.DecodeString(exclusiveStartKey)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid next_key format"})
				return
			}
			json.Unmarshal(keyJSON, &lastEvaluatedKey)
		}

		feedRepo := repository.NewDynamoDBFeedRepository(dynamoDBClient)
		postRepo := repository.NewDynamoDBPostRepository(dynamoDBClient)

		// 1. 從 UserFeed 表獲取 PostID 列表
		paginatedFeed, err := feedRepo.GetUserFeed(c.Request.Context(), userID, int32(limit), lastEvaluatedKey) //
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user feed"})
			return
		}

		if len(paginatedFeed.Items) == 0 {
			c.JSON(http.StatusOK, gin.H{"data": []interface{}{}, "next_key": ""})
			return
		}

		// 2. 提取 PostID 並批量獲取完整貼文內容
		var postIDs []string
		for _, item := range paginatedFeed.Items {
			postIDs = append(postIDs, item.PostID) //
		}

		posts, err := postRepo.GetPostsByIDs(c.Request.Context(), postIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch full posts for feed"})
			return
		}

		likedStatusMap, err := postRepo.CheckIfPostsLikedBy(c.Request.Context(), postIDs, viewerID)
		if err != nil {
			// 即使檢查失敗，我們仍然回傳貼文列表，只是 isLiked 可能不準確
			log.Printf("Could not check liked status for viewer %s: %v", viewerID, err)
			likedStatusMap = make(map[string]bool) // 建立一個空的 map 以免下方出錯
		}

		// --- 3. 將 []models.Post 轉換為前端所需的 []models.FrontendFeedPost ---
		authorCache := make(map[string]string) // 用於快取作者名稱，避免重複查詢
		var frontendPosts []models.PostFeedDTO

		for _, post := range posts {
			authorName, found := authorCache[post.AuthorID]
			if !found {
				authorIDUint, _ := strconv.ParseUint(post.AuthorID, 10, 64)
				user, userErr := userRepo.GetUserByID(uint(authorIDUint))
				if userErr != nil {
					log.Printf("Could not fetch author name for ID %s: %v", post.AuthorID, userErr)
					authorName = "未知的使用者" // 設定備用名稱
				} else {
					authorName = user.Username
				}
				authorCache[post.AuthorID] = authorName
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
				IsLiked:      likedStatusMap[post.PostID],
			}
			frontendPosts = append(frontendPosts, dto)
		}
		
		// 4. (重要) 重新排序，以符合原始 Feed 的時間順序
		postOrder := make(map[string]int)
		for i, item := range paginatedFeed.Items {
			postOrder[item.PostID] = i
		}
		sort.SliceStable(frontendPosts, func(i, j int) bool {
			return postOrder[frontendPosts[i].PostID] < postOrder[frontendPosts[j].PostID]
		})


		// 5. 準備下一次分頁的 next_key
		var nextKey string
		if paginatedFeed.LastEvaluatedKey != nil {
			keyJSON, _ := json.Marshal(paginatedFeed.LastEvaluatedKey)
			nextKey = base64.StdEncoding.EncodeToString(keyJSON)
		}

		c.JSON(http.StatusOK, gin.H{"data": frontendPosts, "next_key": nextKey})
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