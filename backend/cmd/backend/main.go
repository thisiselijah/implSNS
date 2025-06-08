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
	"backend/internal/middleware"
)

func main() {
	// ... 其他初始化程式碼 ...
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Fail to load configurations: %v", err)
	}

	mysqlDB, err := db.InitMySQL(
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)
	if err != nil {
		log.Fatalf("Fail to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	awsdynamoDB, err := db.InitDynamoDB(
		cfg.DynamoDB.Region,
		cfg.DynamoDB.Endpoint,
		cfg.DynamoDB.AccessKey,
		cfg.DynamoDB.SecretKey,
		cfg.DynamoDB.SessionToken,
	)
	if err != nil {
		log.Fatalf("Fail to connect to DynamoDB: %v", err)
	}
	// Repositories
	userRepo := repository.NewMySQLUserRepository(mysqlDB)
	tokenBlacklistRepo := repository.NewMemoryTokenBlacklist()
	postRepo := repository.NewDynamoDBPostRepository(awsdynamoDB)
	// Services
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, cfg.JWT.SecretKey, cfg.JWT.ExpiryMinutes)
	profileService := service.NewProfileService(userRepo)
	postService := service.NewPostService(postRepo, userRepo) 
	userService := service.NewUserService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(*authService, cfg.JWT.ExpiryMinutes)
	profileHandler := handler.NewProfileHandler(profileService)
	postHandler := handler.NewPostHandler(postService)
	userHandler := handler.NewUserHandler(userService, mysqlDB, awsdynamoDB) 

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenBlacklistRepo, cfg.JWT.SecretKey)


	// 6. 初始化 Router
	r := router.NewRouter(mysqlDB, awsdynamoDB, authHandler, profileHandler, postHandler, userHandler, userRepo, authMiddleware)


	// 7. 啟動伺服器
	log.Println("Server starting on port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}