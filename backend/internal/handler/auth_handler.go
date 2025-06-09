package handler

import (

	"net/http"
	"strings"
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
	cookieValue := loginResponse.UserID // 假設 UserID 是字串類型
	maxAgeSeconds := h.jwtTokenExpiryMinutes * 60

	secureCookie := false

	cookie := &http.Cookie{
		Name:     "user_id",
		Value:    cookieValue,
		MaxAge:   maxAgeSeconds,
		Path:     "/",
		Domain:   "",
		Secure:   secureCookie,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)

	// 新增 JWT token 的 cookie
	jwtCookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    loginResponse.Token, // 假設 Token 在這裡
		MaxAge:   maxAgeSeconds,
		Path:     "/",
		Domain:   "",
		Secure:   secureCookie,
		HttpOnly: true, // 建議設為 true
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, jwtCookie)
	// --- Cookie 設定完畢 ---

	// return JSON 格式的登入回應
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"userID":   loginResponse.UserID,       // 使用者 ID
	})
}

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

func (h *AuthHandler) Logout(c *gin.Context) {
    // 從 cookie 讀取 jwt_token
    tokenString, err := c.Cookie("jwt_token")
    if err != nil || tokenString == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_token cookie is required"})
        return
    }

    err = h.authService.Logout(tokenString)
    if err != nil {
        if strings.Contains(err.Error(), "invalid token") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed: " + err.Error()})
        }
        return
    }

    // 清除 jwt_token 與 user_id cookie
    clearCookie := func(name string) {
        http.SetCookie(c.Writer, &http.Cookie{
            Name:     name,
            Value:    "",
            Path:     "/",
            MaxAge:   -1,
            HttpOnly: true,
            SameSite: http.SameSiteLaxMode,
        })
    }
    clearCookie("jwt_token")
    clearCookie("user_id")

    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) GetAuthStatus(c *gin.Context) {
    // 從 cookie 讀取 jwt_token
    tokenString, err := c.Cookie("jwt_token")
    if err != nil || tokenString == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt_token cookie is required"})
        return
    }

    // 檢查 token 是否有效，並取得 userID
    userID, err := h.authService.GetUserIDFromToken(tokenString)
    if err != nil || userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "User is authenticated",
        "userID":  userID,
    })
}