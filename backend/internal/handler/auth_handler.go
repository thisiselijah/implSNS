package handler

import (
	"fmt" // 用於格式化 user ID 為字串
	"net/http"
	"strings"
	// "time" // 如果需要在這裡直接計算 Expires，但我們用 MaxAge

	"backend/internal/models"
	"backend/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 結構體持有 AuthService 的依賴以及 JWT 過期時間（分鐘）
type AuthHandler struct {
	authService           service.AuthService
	jwtTokenExpiryMinutes int // 新增此欄位
}

// NewAuthHandler 是 AuthHandler 的建構子，增加 jwtTokenExpiryMinutes 參數
func NewAuthHandler(authService service.AuthService, jwtTokenExpiryMinutes int) *AuthHandler {
	return &AuthHandler{
		authService:           authService,
		jwtTokenExpiryMinutes: jwtTokenExpiryMinutes, // 儲存 JWT 過期分鐘數
	}
}

// RegisterPayload 定義了註冊請求預期的 JSON 結構
type RegisterPayload struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginPayload 定義了登入請求預期的 JSON 結構
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login 處理登入邏輯，接收 JSON 格式的使用者憑證
func (h *AuthHandler) Login(c *gin.Context) {
	var payload LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	loginData := models.UserForLogin{
		Email:    payload.Email,
		Password: payload.Password,
	}

	loginResponse, err := h.authService.Login(loginData) //
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed: " + err.Error()})
		}
		return
	}

	// --- 設定 HTTP-only cookie ---
	cookieName := "user_id"                                    // 您可以自訂 cookie 名稱
	cookieValue := fmt.Sprintf("%d", loginResponse.UserID)     // Cookie 的值是 UserID
	maxAgeSeconds := h.jwtTokenExpiryMinutes * 60              // Cookie 過期時間（秒），與 JWT 同步

	// 決定 Secure 屬性：
	// 在生產環境 (Release Mode) 且使用 HTTPS 時應為 true。
	// 為簡化本地 HTTP 開發，這裡先設為 false。
	// isReleaseMode := gin.Mode() == gin.ReleaseMode
	secureCookie := false
	// if isReleaseMode { // 如果在生產環境並且你確定是 HTTPS
	// 	secureCookie = true
	// }

	// 設定 Cookie
	// 使用 http.Cookie 結構體可以更完整地設定 SameSite 等屬性
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		MaxAge:   maxAgeSeconds,
		Path:     "/",    // Cookie 在整個網站根路徑下有效
		Domain:   "",     // Domain 留空，瀏覽器會使用當前請求的主機。
		Secure:   secureCookie, // 若為 true，則只在 HTTPS 下傳輸
		HttpOnly: true,         // 核心！設置為 HTTP-only，前端 JS 無法讀取
		SameSite: http.SameSiteLaxMode, // 建議的 SameSite 策略，有助於防止 CSRF
	}
	http.SetCookie(c.Writer, cookie)
	// --- Cookie 設定完畢 ---

	// 登入成功，返回 service 提供的 LoginResponse (包含 token 等)
	c.JSON(http.StatusOK, loginResponse)
}

// Register 處理註冊邏輯 (保持不變)
func (h *AuthHandler) Register(c *gin.Context) {
	var payload RegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	registrationData := models.UserForRegistration{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	registeredUser, err := h.authService.Register(registrationData) //
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if strings.Contains(err.Error(), "password must be at least") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Registration successful",
		"userID":   registeredUser.ID,       //
		"username": registeredUser.Username, //
		"email":    registeredUser.Email,    //
	})
}

// Logout 處理登出邏輯 (保持不變)
func (h *AuthHandler) Logout(c *gin.Context) {
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

	err := h.authService.Logout(tokenString) //
	if err != nil {
		if strings.Contains(err.Error(), "invalid token") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}