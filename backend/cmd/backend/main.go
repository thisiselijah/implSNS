package main

import (
	"log"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handler"   // 引入 handler 套件
	"backend/internal/repository" // 引入 repository 套件
	"backend/internal/router"
	"backend/internal/service"   // 引入 service 套件

	// "github.com/gin-gonic/gin" // 如果 NewRouter 直接返回 *gin.Engine，可能不需要在此 import
	// "database/sql" // db.InitMySQL 返回 *sql.DB
	// "github.com/aws/aws-sdk-go/service/dynamodb" // db.InitDynamoDB 返回 *dynamodb.DynamoDB
)

func main() {
	// 1. 載入設定檔
	cfg, err := config.LoadConfig("config/config.yaml") //
	if err != nil { //
		log.Fatalf("Fail to load configurations: %v", err) //
	}

	// 2. 初始化資料庫連線
	mysqlDB, err := db.InitMySQL( //
		cfg.Database.Username, //
		cfg.Database.Password, //
		cfg.Database.Host,     //
		cfg.Database.Name,     //
	)
	if err != nil { //
		log.Fatalf("Fail to connect to MySQL: %v", err) //
	}

	// DynamoDB 初始化 (如果其他服務會用到)
	awsdynamoDB, err := db.InitDynamoDB( //
		cfg.DynamoDB.Region,       //
		cfg.DynamoDB.Endpoint,     //
		cfg.DynamoDB.AccessKey,    //
		cfg.DynamoDB.SecretKey,    //
		cfg.DynamoDB.SessionToken, //
	)
	if err != nil { //
		log.Fatalf("Fail to connect to DynamoDB: %v", err) //
	}

	// 3. 初始化 Repositories
	// 假設 NewMySQLUserRepository 位於 internal/repository 套件
	// 並且返回一個實現了 repository.UserRepository 介面的實例
	userRepo := repository.NewMySQLUserRepository(mysqlDB)

	// 使用上面定義的簡易 memoryTokenBlacklist 作為 TokenBlacklistRepository 的實現
	// 理想情況下，NewMemoryTokenBlacklist 也應在 repository 套件中
	tokenBlacklistRepo := repository.NewMemoryTokenBlacklist()


	// 4. 初始化 AuthService
	// 確保 cfg.JWT.SecretKey 和 cfg.JWT.ExpiryMinutes 已在 config 中定義並載入
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, cfg.JWT.SecretKey, cfg.JWT.ExpiryMinutes)

	// 5. 初始化 AuthHandler
	authHandler := handler.NewAuthHandler(*authService)


	// 6. 初始化 Router，並傳入 AuthHandler
	// 你需要確保 router.NewRouter 函數的簽名已更新，以接收 authHandler
	// 例如: func NewRouter(mysqlDB *sql.DB, dynamoDBClient *dynamodb.DynamoDB, authHdlr *handler.AuthHandler) *gin.Engine
	r := router.NewRouter(mysqlDB, awsdynamoDB, authHandler) // // 假設 NewRouter 已修改

	// 7. 啟動伺服器
	log.Println("Server starting on port :8080") //
	if err := r.Run(":8080"); err != nil { //
		log.Fatalf("Failed to run server: %v", err) //
	}
}