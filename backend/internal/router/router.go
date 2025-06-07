// backend/internal/router/router.go
package router

import (
	"database/sql"
	"time" // 需要引入 time 套件來設定 MaxAge

	"backend/internal/handler"
	"backend/internal/repository" // <--- 新增導入 repository
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter 負責初始化 Gin 引擎並設定所有應用程式的路由
// 新增 userRepo repository.UserRepository 參數
func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.Client, authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, postHandler *handler.PostHandler, userRepo repository.UserRepository) *gin.Engine { // <--- 修改簽名
	r := gin.Default()

	// --- 設定 CORS 中介軟體 ---
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.2.13:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// --- API 路由 ---
	apiV1 := r.Group("/api/v1")

	// 認證路由
	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	testRoutes := apiV1.Group("/test")
	{
		testRoutes.GET("/mysql/tables", handler.GetTables(mysqlDB))
		testRoutes.GET("/dynamodb//tables", handler.GetDynamoDBTables(dynamoDBClient))
	}

	userRoutes := apiV1.Group("/users")
	{
		userRoutes.GET("/search", )    // 搜尋使用者
		userRoutes.POST("/follow", )   // 追蹤使用者
		userRoutes.POST("/unfollow", ) // 取消追蹤使用者
		userRoutes.GET("/followers/:userID", ) // 獲取使用者的追蹤者
		userRoutes.GET("/following/:userID", ) // 獲取使用者追蹤的使用者

	}

	pagesRoutes := apiV1.Group("/pages")
	{
		// --- Posts ---
		// 將 postHandler 的方法綁定到路由
		pagesRoutes.POST("/posts", postHandler.CreatePost)               // Create a new posts
		pagesRoutes.GET("/posts/:userID", postHandler.GetPostsByUserID) // Get posts by author
		pagesRoutes.POST("/posts/delete", postHandler.DeletePost)         // Delete a post
		pagesRoutes.PUT("/posts/edit", postHandler.UpdatePost)             // Edit a post

		pagesRoutes.POST("/posts/like/:postID", )    // Like a post
		pagesRoutes.POST("/posts/unlike/:postID", )  // Unlike a post
		pagesRoutes.POST("/posts/comment/:postID", ) // Comment on a post
		pagesRoutes.POST("/posts/delete-comment/:postID", ) // Delete a comment on a post

		// Feed
		pagesRoutes.GET("/posts/feed/:userID", handler.GetFeedPosts(dynamoDBClient, userRepo))

		// Profile
		profileRoutes := pagesRoutes.Group("/profile/:userID") // 將 userID 作為共同前綴
		{
			profileRoutes.GET("", profileHandler.GetProfileByUserID)
			profileRoutes.PUT("/avatar", profileHandler.UpdateAvatar)
			profileRoutes.PUT("/bio", profileHandler.UpdateBio)
		}
	}

	return r
}