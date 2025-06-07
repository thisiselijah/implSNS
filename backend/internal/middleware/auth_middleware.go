package middleware

import (
	"backend/internal/service"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware 負責保存驗證中介軟體所需的依賴。
type AuthMiddleware struct {
	blacklistRepo service.TokenBlacklistRepository
	jwtSecretKey  []byte
}

// NewAuthMiddleware 建立一個新的 AuthMiddleware 實例。
func NewAuthMiddleware(blacklistRepo service.TokenBlacklistRepository, secretKey string) *AuthMiddleware {
	return &AuthMiddleware{
		blacklistRepo: blacklistRepo,
		jwtSecretKey:  []byte(secretKey),
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 從 cookie 獲取 jwt_token
        tokenString, err := c.Cookie("jwt_token")
        if err != nil || tokenString == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "jwt_token cookie is required"})
            return
        }

        // 2. 檢查 token 是否在黑名單中
        isBlacklisted, err := m.blacklistRepo.IsTokenBlacklisted(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error checking token status"})
            return
        }
        if isBlacklisted {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated (logged out)"})
            return
        }

        // 3. 解析並驗證 token
        claims := &jwt.RegisteredClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return m.jwtSecretKey, nil
        })

        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
                return
            }
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            return
        }

        if !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

		log.Printf("Authenticated user ID: %s", claims.Subject) 

        // 4. 將 user ID 設定到 context 中，供後續 handler 使用
        c.Set("userID", claims.Subject)

        // 5. 繼續處理下一個 handler
        c.Next()
    }
}
