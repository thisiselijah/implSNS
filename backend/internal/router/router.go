package router

import (
	"database/sql" // (來自之前的 router.go 範例)

	"backend/internal/handler" // (來自之前的 router.go 範例)
	"github.com/aws/aws-sdk-go-v2/service/dynamodb" // (來自之前的 router.go 範例)
	"github.com/gin-gonic/gin" // (來自之前的 router.go 範例)
)

// NewRouter 負責初始化 Gin 引擎並設定所有應用程式的路由
func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.Client, authHandler *handler.AuthHandler) *gin.Engine { // 新增 authHandler 參數
	r := gin.Default() // (來自之前的 router.go 範例)

	// --- API 路由 ---
	apiV1 := r.Group("/api/v1") // 範例分組

	// 認證路由
	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		// Logout 通常需要認證中間件來保護，確保只有登入的用戶才能登出
		// 例如: authRequired := AuthMiddleware() // 你需要自己實現這個中間件
		// authRoutes.POST("/logout", authRequired, authHandler.Logout)
		authRoutes.POST("/logout", authHandler.Logout) // 簡化版，未加中間件
	}

	// 其他 MySQL 相關的路由
	// 假設 handler.GetTables 仍然存在且不需要 authHandler
	otherMySQLRoutes := apiV1.Group("/mysql")
	{
		otherMySQLRoutes.GET("/tables", handler.GetTables(mysqlDB)) // (來自之前的 router.go 範例)
	}

	// DynamoDB 相關的路由
	// 假設 handler.GetDynamoDBTables 仍然存在
	otherDynamoDBRoutes := apiV1.Group("/dynamodb")
	{
		otherDynamoDBRoutes.GET("/tables", handler.GetDynamoDBTables(dynamoDBClient)) // (來自之前的 router.go 範例)
	}


	// --- 健康檢查或其他通用路由 ---
	// r.GET("/health", handler.HealthCheck) // (來自之前的 router.go 範例)

	return r // (來自之前的 router.go 範例)
}