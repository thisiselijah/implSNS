package models

import (
	"database/sql" // 需要引入，因為 BirthDate 可能是 sql.NullTime
	"time"
)

// UserProfile 對應資料庫中的 user_profiles 表
type UserProfile struct {
	ID        uint         `json:"id"`
	UserID    uint         `json:"user_id"`
	Username  string       `json:"username" binding:"required,max=255"` // 限制最大長度為 50 字元
	AvatarURL string       `json:"avatar_url"`
	BirthDate sql.NullTime `json:"birth_date"` // 使用 sql.NullTime 來處理可能的 NULL 值
	Bio       string       `json:"bio"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// UpdateBioPayload 定義了更新個人簡介時請求的 JSON 結構
type UpdateBioPayload struct {
	Bio string `json:"bio" binding:"required"`
}

// UpdateAvatarPayload 定義了更新頭像時請求的 JSON 結構
type UpdateAvatarPayload struct {
	AvatarURL string `json:"avatar_url" binding:"required,url"` // 確保 AvatarURL 是有效的 URL
}