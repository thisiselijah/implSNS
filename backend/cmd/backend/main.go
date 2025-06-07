// backend/cmd/backend/main.go
package main

import (
	"log"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/router"
	"backend/internal/service"
)

func main() {
	// 1. 載入設定檔
	cfg, err := config.LoadConfig("config/config.yaml") //
	if err != nil {
		log.Fatalf("Fail to load configurations: %v", err) //
	}

	// 2. 初始化資料庫連線
	mysqlDB, err := db.InitMySQL( //
		cfg.Database.Username, //
		cfg.Database.Password, //
		cfg.Database.Host,     //
		cfg.Database.Name,     //
	)
	if err != nil {
		log.Fatalf("Fail to connect to MySQL: %v", err) //
	}
	defer mysqlDB.Close()

	awsdynamoDB, err := db.InitDynamoDB( //
		cfg.DynamoDB.Region,       //
		cfg.DynamoDB.Endpoint,     //
		cfg.DynamoDB.AccessKey,    //
		cfg.DynamoDB.SecretKey,    //
		cfg.DynamoDB.SessionToken, //
	)
	if err != nil {
		log.Fatalf("Fail to connect to DynamoDB: %v", err) //
	}

	// 3. 初始化 Repositories
	userRepo := repository.NewMySQLUserRepository(mysqlDB)
	tokenBlacklistRepo := repository.NewMemoryTokenBlacklist()
	postRepo := repository.NewDynamoDBPostRepository(awsdynamoDB) // <-- 新增 PostRepository

	// 4. 初始化 Services
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, cfg.JWT.SecretKey, cfg.JWT.ExpiryMinutes)
	profileService := service.NewProfileService(userRepo)
	postService := service.NewPostService(postRepo) // <-- 新增 PostService

	// 5. 初始化 Handlers
	authHandler := handler.NewAuthHandler(*authService, cfg.JWT.ExpiryMinutes)
	profileHandler := handler.NewProfileHandler(profileService)
	postHandler := handler.NewPostHandler(postService) // <-- 新增 PostHandler

	// 6. 初始化 Router
	r := router.NewRouter(mysqlDB, awsdynamoDB, authHandler, profileHandler, postHandler, userRepo) // <-- 傳入 postHandler

	// 7. 啟動伺服器
	log.Println("Server starting on port :8080") //
	if err := r.Run(":8080"); err != nil { //
		log.Fatalf("Failed to run server: %v", err) //
	}
}