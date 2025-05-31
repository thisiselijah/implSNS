package handler

import (
	"net/http" // 引入 net/http 以便使用 http 狀態碼常數
	"strings"  // 用於處理字串，例如提取 Bearer token

	"backend/internal/models"  // 引入 Service 層期望的 DTO
	"backend/internal/service" // 引入 AuthService 介面或實例
	"github.com/gin-gonic/gin"
)

// AuthHandler 結構體持有 AuthService 的依賴
type AuthHandler struct {
	authService service.AuthService // 注意：這裡最好是 service.AuthService 介面
}

// NewAuthHandler 是 AuthHandler 的建構子
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterPayload 定義了註冊請求預期的 JSON 結構
type RegisterPayload struct {
	Username string `json:"username" binding:"required"`      //
	Email    string `json:"email" binding:"required,email"` //
	Password string `json:"password" binding:"required,min=8"`  //
}

// LoginPayload 定義了登入請求預期的 JSON 結構
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"` //
	Password string `json:"password" binding:"required"`    //
}

// Login 處理登入邏輯，接收 JSON 格式的使用者憑證
func (h *AuthHandler) Login(c *gin.Context) { // 改為 AuthHandler 的方法
	var payload LoginPayload //

	if err := c.ShouldBindJSON(&payload); err != nil { //
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()}) //
		return
	}

	// 將 handler 的 LoginPayload 轉換為 service 層的 models.UserForLogin
	loginData := models.UserForLogin{
		Email:    payload.Email,
		Password: payload.Password,
	}

	// 呼叫 AuthService 的 Login 方法
	loginResponse, err := h.authService.Login(loginData)
	if err != nil {
		// 根據 service 返回的錯誤類型決定 HTTP 狀態碼
		// 例如，"invalid email or password" 通常是 401 Unauthorized
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			// 其他內部錯誤
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		}
		return
	}

	// 登入成功，返回 service 提供的 LoginResponse (包含 token 等)
	c.JSON(http.StatusOK, loginResponse)
}

// Register 處理註冊邏輯，接收 JSON 格式的使用者註冊資訊
func (h *AuthHandler) Register(c *gin.Context) { // 改為 AuthHandler 的方法
	var payload RegisterPayload //

	if err := c.ShouldBindJSON(&payload); err != nil { //
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()}) //
		return
	}

	// 將 handler 的 RegisterPayload 轉換為 service 層的 models.UserForRegistration
	registrationData := models.UserForRegistration{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	// 呼叫 AuthService 的 Register 方法
	// AuthService 的 Register 方法返回 (*models.User, error)
	registeredUser, err := h.authService.Register(registrationData)
	if err != nil {
		// 根據 service 返回的錯誤類型決定 HTTP 狀態碼
		// 例如 "username already exists" 或 "email already exists" 通常是 409 Conflict
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "password must be at least") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// 其他內部錯誤
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		}
		return
	}

	// 註冊成功，通常返回 201 Created
	// 你可以選擇返回部分用戶資訊（不含敏感資訊）
	c.JSON(http.StatusCreated, gin.H{ //
		"message":  "Registration successful", //
		"userID":   registeredUser.ID,
		"username": registeredUser.Username, //
		"email":    registeredUser.Email,    //
	})
}

// Logout 處理登出邏輯
func (h *AuthHandler) Logout(c *gin.Context) { // 改為 AuthHandler 的方法
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return
	}
	tokenString := parts[1]

	// 呼叫 AuthService 的 Logout 方法
	err := h.authService.Logout(tokenString)
	if err != nil {
		if strings.Contains(err.Error(), "invalid token") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"}) //
}