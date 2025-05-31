package repository

import (
	"time"
	"log" // 用於簡單的日誌輸出
)

// TokenBlacklistRepository 定義了 token 黑名單的操作
type TokenBlacklistRepository interface {
	BlacklistToken(tokenString string, expiresAt time.Time) error
	IsTokenBlacklisted(tokenString string) (bool, error)
}
// --- MemoryTokenBlacklist (簡易記憶體黑名單實現範例) ---
// 在實際應用中，你可能會使用 Redis 或其他持久化儲存。
// 這個結構和其方法應該放在 repository 套件中，例如 internal/repository/memory_token_repository.go
type memoryTokenBlacklist struct {
	list map[string]time.Time
}

// NewMemoryTokenBlacklist 建立一個基於記憶體的 TokenBlacklistRepository 實例
func NewMemoryTokenBlacklist() TokenBlacklistRepository { // 確保返回的是定義在 repository 套件的介面
	return &memoryTokenBlacklist{list: make(map[string]time.Time)}
}

// BlacklistToken 將 token 加入黑名單
func (m *memoryTokenBlacklist) BlacklistToken(tokenString string, expiresAt time.Time) error {
	m.list[tokenString] = expiresAt
	log.Printf("Token blacklisted in memory: %s (expires at %v)", tokenString, expiresAt)
	// 實際應用中，可能需要一個 goroutine 定期清理已過期的 token
	return nil
}

// IsTokenBlacklisted 檢查 token 是否在黑名單中且未過期
func (m *memoryTokenBlacklist) IsTokenBlacklisted(tokenString string) (bool, error) {
	expiry, found := m.list[tokenString]
	if found && time.Now().Before(expiry) {
		return true, nil // 存在且未過期
	}
	// 如果 token 存在但已過期，可以從黑名單中移除 (可選)
	if found && time.Now().After(expiry) {
		delete(m.list, tokenString)
	}
	return false, nil
}
// --- END MemoryTokenBlacklist ---