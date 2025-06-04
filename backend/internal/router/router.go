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
func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.Client, authHandler *handler.AuthHandler, userRepo repository.UserRepository) *gin.Engine { // <--- 修改簽名
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
	
	postRoutes := apiV1.Group("/pages")
	{
		// Feed 
		// 將 userRepo 傳遞給 GetFeedPosts
		postRoutes.GET("/posts/feed/:userID", handler.GetFeedPosts(dynamoDBClient, userRepo)) // <--- 修改呼叫
	}

	return r
}