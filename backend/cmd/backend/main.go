// backend/cmd/backend/main.go
// GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/backend/main ./cmd/backend/
// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/backend/

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
	"backend/internal/recommendation"
	"context"
	"time"
)

// startTrendingRecommendationGenerator 在背景定期執行推薦演算法
func startTrendingRecommendationGenerator(recommender *recommendation.TrendingRecommender, interval time.Duration) {
	log.Printf("Starting periodic trending recommendation generator with interval %v", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 為了讓服務啟動後能立即提供推薦，先執行一次
	log.Println("Running initial trending recommendation generation on startup...")
	if err := recommender.GenerateRecommendations(context.Background()); err != nil {
		log.Printf("Error during initial recommendation generation: %v", err)
	}


	// 進入定期執行的循環
	for {
		select {
		case <-ticker.C:
			log.Println("Ticker fired. Starting scheduled trending recommendation generation...")
			// 每次執行都使用一個新的背景 context
			err := recommender.GenerateRecommendations(context.Background())
			if err != nil {
				log.Printf("Error during scheduled recommendation generation: %v", err)
			} else {
				log.Println("Scheduled trending recommendation generation finished successfully.")
			}
		}
	}
}

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
	feedRepo := repository.NewDynamoDBFeedRepository(awsdynamoDB)
	recoRepo := repository.NewDynamoDBRecommendationRepository(awsdynamoDB)

	// Recommendation 

	trendingRecommender := recommendation.NewTrendingRecommender(postRepo, userRepo, recoRepo) //

	// --- 啟動背景任務 ---
	// 使用 goroutine 執行，才不會阻塞主線程的 Web 伺服器啟動
	go startTrendingRecommendationGenerator(trendingRecommender, 1*time.Hour) //


	// Services
	authService := service.NewAuthService(userRepo, tokenBlacklistRepo, cfg.JWT.SecretKey, cfg.JWT.ExpiryMinutes)
	profileService := service.NewProfileService(userRepo)
	postService := service.NewPostService(postRepo, userRepo, feedRepo) 
	userService := service.NewUserService(userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(*authService, cfg.JWT.ExpiryMinutes)
	profileHandler := handler.NewProfileHandler(profileService)
	postHandler := handler.NewPostHandler(postService, userRepo, feedRepo, postRepo, recoRepo)
	userHandler := handler.NewUserHandler(userService, mysqlDB, awsdynamoDB) 

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenBlacklistRepo, cfg.JWT.SecretKey)


	// 6. 初始化 Router
	r := router.NewRouter(mysqlDB, awsdynamoDB, authHandler, profileHandler, postHandler, userHandler, userRepo, authMiddleware) //



	// 7. 啟動伺服器
	log.Println("Server starting on port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}