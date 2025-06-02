package router

import (
	"database/sql"
	"time" // 需要引入 time 套件來設定 MaxAge

	"backend/internal/handler"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-contrib/cors" // <--- 引入 CORS 中介軟體
	"github.com/gin-gonic/gin"
)

// NewRouter 負責初始化 Gin 引擎並設定所有應用程式的路由
func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.Client, authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()

	// --- 設定 CORS 中介軟體 ---
	// 這是一個建議的基礎設定，你可以根據你的需求調整
	config := cors.Config{
		// AllowOrigins 指定了允許哪些來源 (前端的 URL) 進行跨來源請求。
		// 你前端的 auth.js 顯示來源是 http://localhost:3000
		// 如果你的前端也可能透過 IP 存取，例如 http://192.168.0.144:3000 (假設前端和後端在同一台機器但埠號不同)，也可以加入。
		AllowOrigins: []string{"http://localhost:3000"}, // 根據你的前端實際來源調整

		// AllowMethods 指定了允許的 HTTP 方法。
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},

		// AllowHeaders 指定了允許的請求標頭。
		// "Authorization" 是 JWT 常用來傳遞 token 的標頭。
		// "Content-Type" 是你的 auth.js 中有設定的。
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},

		// ExposeHeaders 允許前端 JavaScript 存取的回應標頭 (預設情況下，前端只能存取部分簡單標頭)。
		// 例如，如果你的後端在回應中設定了自訂標頭，且希望前端能讀取，就需要在此列出。
		ExposeHeaders: []string{"Content-Length"},

		// AllowCredentials 指示是否允許請求攜帶憑證 (如 cookies 或 HTTP 認證)。
		// 如果你的前端需要傳送 cookies，或者將來會用到 session，可以設為 true。
		// 對於 Bearer token 認證，這個設定通常不是嚴格必要的，但包含它通常無害。
		AllowCredentials: true,

		// MaxAge 指示預檢請求 (OPTIONS request) 的結果可以被瀏覽器快取多久 (單位是秒)。
		MaxAge: 12 * time.Hour,

		// AllowAllOrigins: true, // 如果你想允許任何來源 (等同於 Access-Control-Allow-Origin: *)
		// 但這通常不建議在生產環境中使用，因為安全性較低。明確指定 AllowOrigins 比較好。
	}
	r.Use(cors.New(config)) // <--- 將 CORS 中介軟體應用到 Gin 引擎

	// --- API 路由 ---
	apiV1 := r.Group("/api/v1") //

	// 認證路由
	authRoutes := apiV1.Group("/auth") //
	{
		authRoutes.POST("/register", authHandler.Register) //
		authRoutes.POST("/login", authHandler.Login)       //
		authRoutes.POST("/logout", authHandler.Logout)     //
	}

	// 其他 MySQL 相關的路由
	otherMySQLRoutes := apiV1.Group("/mysql") //
	{
		otherMySQLRoutes.GET("/tables", handler.GetTables(mysqlDB)) //
	}

	// DynamoDB 相關的路由
	otherDynamoDBRoutes := apiV1.Group("/dynamodb") //
	{
		otherDynamoDBRoutes.GET("/tables", handler.GetDynamoDBTables(dynamoDBClient)) //
	}

	return r
}