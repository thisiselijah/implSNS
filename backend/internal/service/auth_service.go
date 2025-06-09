package service

import (
	"errors" // 用於建立自訂錯誤
	"log"    // 簡單日誌記錄
	"time"   // 用於 JWT 的過期時間
	"backend/internal/models"
	"backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)


// --- Token Blacklist Interface (模擬，用於登出) ---
type TokenBlacklistRepository interface {
	BlacklistToken(tokenString string, expiresAt time.Time) error
	IsTokenBlacklisted(tokenString string) (bool, error)
}

// --- AuthService ---
type AuthService struct {
	userRepo       repository.UserRepository
	blacklistRepo  TokenBlacklistRepository // 用於登出時將 token 加入黑名單
	jwtSecretKey   []byte                 // JWT 簽名用的密鑰
	jwtTokenExpiry time.Duration          // JWT 過期時間
}

// NewAuthService 是 AuthService 的建構子
func NewAuthService(userRepo repository.UserRepository, blacklistRepo TokenBlacklistRepository, secretKey string, tokenExpiryMinutes int) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		blacklistRepo:  blacklistRepo,
		jwtSecretKey:   []byte(secretKey),
		jwtTokenExpiry: time.Minute * time.Duration(tokenExpiryMinutes),
	}
}

// Register 處理使用者註冊邏輯
func (s *AuthService) Register(userData models.UserForRegistration) (*models.User, error) {
	// 1. 檢查使用者名稱是否已存在
	if _, err := s.userRepo.GetUserByUsername(userData.Username); err == nil {
		// 如果 err 為 nil，表示找到了使用者，因此使用者名稱已存在
		return nil, errors.New("username already exists")
	}

	// 2. 檢查 Email 是否已存在
	if _, err := s.userRepo.GetUserByEmail(userData.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 3. 密碼強度檢查 (這裡可以加入更複雜的邏輯)
	if len(userData.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// 4. 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("failed to process registration") // 不要回傳內部錯誤細節
	}

	// 5. 建立 User 物件
	now := time.Now()
	newUser := &models.User{
		Username:     userData.Username,
		Email:        userData.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 6. 儲存到資料庫
	if err := s.userRepo.CreateUser(newUser); err != nil {
		log.Printf("Error creating user in DB: %v", err)
		return nil, errors.New("failed to register user")
	}

	// 為了安全，返回的 User 物件不應包含 PasswordHash
	// newUser.PasswordHash = "" // 如果直接返回 newUser，可以清空
	// 或者返回一個不包含敏感資訊的 DTO
	// 在這裡我們直接返回創建的 User 物件 (ID 應該由 DB 填充)
	// 但記得 handler 層不要將 PasswordHash 序列化給前端
	return newUser, nil
}

// Login 處理使用者登入邏輯
func (s *AuthService) Login(loginData models.UserForLogin) (*models.LoginResponse, error) {
	// 1. 根據 Email 查找使用者
	user, err := s.userRepo.GetUserByEmail(loginData.Email)
	if err != nil {
		// 包括使用者不存在或資料庫錯誤的情況
		log.Printf("Login attempt: User with email %s not found or DB error: %v", loginData.Email, err)
		return nil, errors.New("invalid email or password") // 通用錯誤訊息
	}

	// 2. 比對密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password))
	if err != nil {
		// 密碼不匹配
		log.Printf("Login attempt: Password mismatch for email %s", loginData.Email)
		return nil, errors.New("invalid email or password") // 通用錯誤訊息
	}

	// 3. 產生 JWT

	expirationTime := time.Now().Add(s.jwtTokenExpiry)
	claims := &jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "my-app", // 應用程式名稱
		// 你可以添加自訂的 claims
		// "username": user.Username,
		// "email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		log.Printf("Error generating JWT for user %s: %v", loginData.Email, err)
		return nil, errors.New("failed to login, please try again later")
	}

	// 4. 登入成功，返回 Token 和使用者資訊
	return &models.LoginResponse{
		Token:        tokenString,
		UserID:       user.ID,
		UserEmail:    user.Email,
		UserUsername: user.Username,
	}, nil
}

// Logout 處理使用者登出邏輯
// 對於 JWT，一個常見的伺服器端登出策略是將 token 加入黑名單。
// 客戶端也應該在登出時刪除本地儲存的 token。
func (s *AuthService) Logout(tokenString string) error {
	// 1. 解析 token 以獲取其過期時間等資訊 (可選，但有助於黑名單管理)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecretKey, nil
	})

	if err != nil {
		// 如果 token 無效 (例如已過期或簽名不對)，從伺服器角度看，它已經「登出」了
		// 但如果它只是格式不對，我們可能還是想記錄一下
		log.Printf("Logout attempt with invalid token: %v", err)
		return errors.New("invalid token provided for logout")
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// 2. 將 token 加入黑名單，直到它自然過期
		// 這裡假設 BlacklistToken 方法會處理重複加入等情況
		expiresAt := claims.ExpiresAt.Time
		if err := s.blacklistRepo.BlacklistToken(tokenString, expiresAt); err != nil {
			log.Printf("Error blacklisting token: %v", err)
			return errors.New("failed to logout, please try again")
		}
		log.Printf("Token for user (sub: %s) blacklisted until %v", claims.Subject, expiresAt)
		return nil
	}

	return errors.New("invalid token provided for logout")
}

func (s *AuthService) GetUserIDFromToken(tokenString string) (string, error) {
    // 1. 先檢查黑名單
    isBlacklisted, err := s.blacklistRepo.IsTokenBlacklisted(tokenString)
    if err != nil {
        return "", err
    }
    if isBlacklisted {
        return "", errors.New("token is blacklisted")
    }

    // 2. 解析 JWT
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return s.jwtSecretKey, nil
    })
    if err != nil || !token.Valid {
        return "", errors.New("invalid or expired token")
    }

    return claims.Subject, nil // userID 存在於 Subject
}