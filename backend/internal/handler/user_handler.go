package handler

import (
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "context"
)



func GetTables(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        rows, err := db.Query("SHOW TABLES")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取資料表"})
            return
        }
        defer rows.Close()

        var tables []string
        for rows.Next() {
            var table string
            if err := rows.Scan(&table); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取資料表"})
                return
            }
            tables = append(tables, table)
        }

        c.JSON(http.StatusOK, tables)
    }
}

func GetDynamoDBTables(dynamoDBClient *dynamodb.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        input := &dynamodb.ListTablesInput{}

        result, err := dynamoDBClient.ListTables(context.TODO(), input)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取 DynamoDB 資料表"})
            return
        }

        c.JSON(http.StatusOK, result.TableNames)
    }
}

