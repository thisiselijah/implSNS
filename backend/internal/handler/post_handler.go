// internal/handler/post_handler.go
package handler

import (
	"backend/internal/models"      // 引入 models，其中包含 PostFeedDTO 和 Post (from feed_model)
	"backend/internal/repository" // 引入 repository 介面
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv" // 用於將 string 類型的 AuthorID 轉換為 uint
)

// GetFeedPosts 處理獲取使用者 Feed 的請求
// 現在也接收 userRepo (UserRepository) 來查詢作者名稱
func GetFeedPosts(dynamoDBClient *dynamodb.Client, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID") // 從 URL 路徑中獲取 userID

		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
			return
		}

		// 初始化 PostRepository (操作 DynamoDB)
		// 假設 NewDynamoDBPostRepository 已經存在並返回一個實現了所需方法的實例
		postRepo := repository.NewDynamoDBPostRepository(dynamoDBClient)

		// 1. 從 Repository 獲取該 userID 的 Feed Item 列表 (只包含 post_id 等引用資訊)
		feedItems, err := postRepo.GetFeedItemsByUserID(c.Request.Context(), "USER#"+userID)
		if err != nil {
			log.Printf("Error fetching feed items for userID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed items"})
			return
		}

		if len(feedItems) == 0 {
			c.JSON(http.StatusOK, []models.PostFeedDTO{}) // 返回空的 DTO slice
			return
		}

		// 2. 提取所有 post_id
		var postIDs []string
		for _, item := range feedItems {
			if item.PostID != "" {
				postIDs = append(postIDs, item.PostID)
			}
		}

		if len(postIDs) == 0 {
			log.Printf("No valid postIDs found in feed items for userID %s", userID)
			c.JSON(http.StatusOK, []models.PostFeedDTO{}) // 返回空的 DTO slice
			return
		}
		log.Printf("Extracted postIDs for userID %s: %v", userID, postIDs)

		// 3. 根據 post_id 列表，從 Repository 批量獲取完整的貼文內容 (models.Post)
		postsFromDB, err := postRepo.GetPostsByIDs(c.Request.Context(), postIDs)
		if err != nil {
			log.Printf("Error fetching full posts for userID %s, postIDs %v: %v", userID, postIDs, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch full posts for feed"})
			return
		}

		// 4. 將 models.Post 轉換為 models.PostFeedDTO，並填充 AuthorName
		var feedDTOs []models.PostFeedDTO
		// 為了保持 feedItems 的原始順序，我們遍歷 feedItems 並從 postMap 中查找對應的 post
		postMap := make(map[string]models.Post) // models.Post 是 DynamoDB 的原始貼文結構
		for _, p := range postsFromDB {
			postMap[p.PostID] = p
		}

		for _, feedItem := range feedItems {
			if post, ok := postMap[feedItem.PostID]; ok {
				var authorName string
				// models.Post.AuthorID 是 string，需要轉換為 uint 以查詢 userRepo.GetUserByID
				// 假設 DynamoDB 中儲存的 AuthorID 是 User 的數字 ID 的字串形式
				authorIDUint, convErr := strconv.ParseUint(post.AuthorID, 10, 64) // 使用 ParseUint
				if convErr != nil {
					log.Printf("Error converting AuthorID string '%s' to uint for post '%s': %v. Using original AuthorID as name.", post.AuthorID, post.PostID, convErr)
					authorName = "AuthorID: " + post.AuthorID // 降級處理：使用原始 ID
				} else {
					// 透過 UserRepository 查詢作者的使用者名稱
					user, userErr := userRepo.GetUserByID(uint(authorIDUint)) // 轉換為 uint
					if userErr != nil {
						log.Printf("Error fetching user (ID: %d) for post '%s': %v. Using original AuthorID as name.", authorIDUint, post.PostID, userErr)
						authorName = "User: " + post.AuthorID // 降級處理
					} else {
						authorName = user.Username // 假設 models.User 有 Username 欄位
					}
				}

				// 建立 DTO
				dto := models.PostFeedDTO{
					PostID:       post.PostID,       //
					AuthorName:   authorName,        // 使用查詢到的名稱
					Content:      post.Content,      //
					Media:        post.Media,        //
					Tags:         post.Tags,         //
					Location:     post.Location,     //
					LikeCount:    post.LikeCount,    //
					CommentCount: post.CommentCount, //
					CreatedAt:    post.CreatedAt,    //
					UpdatedAt:    post.UpdatedAt,    //
				}
				feedDTOs = append(feedDTOs, dto)
			} else {
				log.Printf("Post with ID %s found in feedItems but not in fetched postsMap for userID %s", feedItem.PostID, userID)
			}
		}

		c.JSON(http.StatusOK, feedDTOs)
	}
}