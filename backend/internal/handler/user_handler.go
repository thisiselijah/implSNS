package handler

import (
    "database/sql"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "context"
    "backend/internal/service"
)
// UserHandler 結構體
type UserHandler struct {
	userService *service.UserService
	db          *sql.DB
	dynamoDBClient *dynamodb.Client
}

// NewUserHandler 是 UserHandler 的建構子
func NewUserHandler(userService *service.UserService, db *sql.DB, dynamo *dynamodb.Client) *UserHandler {
	return &UserHandler{
		userService: userService,
		db:          db,
		dynamoDBClient: dynamo,
	}
}

// GetFollowers 處理獲取粉絲列表的請求
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	followers, err := h.userService.GetFollowers(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// GetFollowing 處理獲取正在追蹤列表的請求
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	following, err := h.userService.GetFollowing(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following list"})
		return
	}

	c.JSON(http.StatusOK, following)
}

// FollowUser 處理追蹤使用者的請求
func (h *UserHandler) FollowUser(c *gin.Context) {
	// 從 URL 參數中獲取要追蹤的使用者 ID
	followedIDStr := c.Param("userID")
	followedID, err := strconv.ParseUint(followedIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 從身份驗證中介軟體設定的 context 中獲取當前登入使用者的 ID
	followerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	followerID, err := strconv.ParseUint(followerIDStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authenticated user ID"})
		return
	}
	
	// 呼叫 service 執行追蹤邏輯
	if err := h.userService.FollowUser(uint(followerID), uint(followedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser 處理取消追蹤使用者的請求
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	// 從 URL 參數中獲取要取消追蹤的使用者 ID
	followedIDStr := c.Param("userID")
	followedID, err := strconv.ParseUint(followedIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// 從 context 中獲取當前登入使用者的 ID
	followerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	followerID, err := strconv.ParseUint(followerIDStr.(string), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authenticated user ID"})
		return
	}

	// 呼叫 service 執行取消追蹤邏輯
	if err := h.userService.UnfollowUser(uint(followerID), uint(followedID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

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

