package models

import (
	"time"
)

// --- User Model (假設的資料庫模型) ---
// 在實際專案中，這個 User 結構體應該在 models 套件中定義
type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // 密碼雜湊不應該被序列化到 JSON
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// --- DTOs (Data Transfer Objects) ---
// 這些結構體可以從 handler 層傳入，或者在 service 層內部轉換得到
// 為了簡潔，我們直接在 service 層定義，實際專案中可能放在 models 或專門的 dto 套件
// UserForRegistration 代表註冊時傳入的資料
type UserForRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserForLogin 代表登入時傳入的資料
type UserForLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse 代表成功登入後的回應
type LoginResponse struct {
	Token        string `json:"token"`
	UserID       uint   `json:"user_id"` // 假設 UserID 是 uint
	UserEmail    string `json:"email"`
	UserUsername string `json:"username"`
}