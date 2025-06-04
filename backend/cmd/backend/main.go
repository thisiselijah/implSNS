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
	userRepo := repository.NewMySQLUserRepository(mysqlDB)          //
	tokenBlacklistRepo := repository.NewMemoryTokenBlacklist()      //

	// 4. 初始化 AuthService
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, cfg.JWT.SecretKey, cfg.JWT.ExpiryMinutes) //

	// 5. 初始化 AuthHandler
	authHandler := handler.NewAuthHandler(*authService, cfg.JWT.ExpiryMinutes) //

	// 6. 初始化 Router, 將 userRepo 傳遞進去
	r := router.NewRouter(mysqlDB, awsdynamoDB, authHandler, userRepo) // <--- 修改此處, 傳入 userRepo

	// 7. 啟動伺服器
	log.Println("Server starting on port :8080") //
	if err := r.Run(":8080"); err != nil { //
		log.Fatalf("Failed to run server: %v", err) //
	}
}