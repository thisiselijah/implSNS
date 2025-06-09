// backend/internal/router/router.go
package router

import (
	"database/sql"
	"time"
	"backend/internal/middleware"
	"backend/internal/handler"
	"backend/internal/repository"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.Client, authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, postHandler *handler.PostHandler, userHandler *handler.UserHandler, userRepo repository.UserRepository, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	// --- CORS 中介軟體設定 ---
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// --- API 路由 ---
	apiV1 := r.Group("/api/v1")

	// --- 公開路由 (無需身份驗證) ---
	authPublicRoutes := apiV1.Group("/auth")
	{
		authPublicRoutes.POST("/register", authHandler.Register)
		authPublicRoutes.POST("/login", authHandler.Login)
	}

    // --- 保護路由 (需要身份驗證) ---
	// 任何使用此中介軟體的路由群組都需要一個有效的 JWT
	authRequired := apiV1.Group("/")
	authRequired.Use(authMiddleware.Authenticate())
	{
		// 登出需要驗證身份，以識別要加入黑名單的 token
		authRequired.POST("/auth/logout", authHandler.Logout)
		authRequired.GET("/auth/status", authHandler.GetAuthStatus)

		// 使用者相關操作
		userRoutes := authRequired.Group("/users")
		{
			userRoutes.POST("/:userID/follow", userHandler.FollowUser)
			userRoutes.POST("/:userID/unfollow", userHandler.UnfollowUser)
			userRoutes.GET("/:userID/followers", userHandler.GetFollowers)
			userRoutes.GET("/:userID/following", userHandler.GetFollowing)
		}

		// 頁面相關內容的群組
		pagesRoutes := authRequired.Group("/pages")
		{
			// --- 貼文 ---
			pagesRoutes.POST("/posts", postHandler.CreatePost)
			pagesRoutes.GET("/posts/:userID", postHandler.GetPostsByUserID)
			pagesRoutes.POST("/posts/delete", postHandler.DeletePost)
			pagesRoutes.PUT("/posts/edit", postHandler.UpdatePost)

			// --- 貼文互動 ---
			postInteractionRoutes := pagesRoutes.Group("/posts/:postID")
			{
				postInteractionRoutes.PUT("/like", postHandler.LikePost)
				postInteractionRoutes.PUT("/unlike", postHandler.UnlikePost)
				postInteractionRoutes.POST("/comment", postHandler.CreateComment)
				postInteractionRoutes.DELETE("/comment/:commentSK", postHandler.DeleteComment)
			}

			// --- 動態消息 (Feed) ---
			pagesRoutes.GET("/posts/feed/:userID", postHandler.GetFeedPosts)

			// --- 個人資料 ---
			profileRoutes := pagesRoutes.Group("/profile/:userID")
			{
				profileRoutes.GET("", profileHandler.GetProfileByUserID)
				profileRoutes.PUT("/avatar", profileHandler.UpdateAvatar)
				profileRoutes.PUT("/bio", profileHandler.UpdateBio)
			}
		}

	}

	return r
}