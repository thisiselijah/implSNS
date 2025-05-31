package main

import (
	"backend/internal/config"
	// "fmt"
	"backend/internal/db"
	"backend/internal/handler"
	"log"
	"github.com/gin-gonic/gin"
)

func main() {
    // fmt.Println("Hello, World!")
    cfg, err := config.LoadConfig("config/config.yaml")
    if err != nil {
        log.Fatalf("Fail to load confugurations.: %v", err)
    }

    mysqlDB, err := db.InitMySQL(
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Name,
    )
    if err != nil {
        log.Fatalf("Fail to connect databse %v", err)
    }

    awsdynamoDB, err := db.InitDynamoDB(
        cfg.DynamoDB.Region,
        cfg.DynamoDB.Endpoint,
        cfg.DynamoDB.AccessKey,
        cfg.DynamoDB.SecretKey,
        cfg.DynamoDB.SessionToken,
    )

    if err != nil {
        log.Fatalf("Fail to connect DynamoDB: %v", err)
    }

    r := gin.Default()
    // MySQL 的路由
    r.GET("/tables", handler.GetTables(mysqlDB))
    
    // DynamoDB 的路由
    r.GET("/dynamodb/tables", handler.GetDynamoDBTables(awsdynamoDB))

    r.Run(":8080")
}

