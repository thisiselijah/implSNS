package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9" // 使用 go-redis/v9
)

const (
	// tokenBlacklistPrefix 是儲存在 Redis 中黑名單 token 的鍵前綴，方便管理和區分
	tokenBlacklistPrefix = "blacklist:token:"
)

// redisTokenBlacklistRepository 實現了 TokenBlacklistRepository 介面
type redisTokenBlacklistRepository struct {
	client *redis.Client // Redis 客戶端
}

// NewRedisTokenBlacklistRepository 是 redisTokenBlacklistRepository 的建構子
func NewRedisTokenBlacklistRepository(client *redis.Client) TokenBlacklistRepository {
	return &redisTokenBlacklistRepository{client: client}
}

// BlacklistToken 將指定的 token 字串加入 Redis 黑名單，直到其 'expiresAt' 時間點。
func (r *redisTokenBlacklistRepository) BlacklistToken(tokenString string, expiresAt time.Time) error {
	ctx := context.Background() // 或者從呼叫者傳入 context

	// 計算 token 剩餘的有效時間
	// Redis 的 SETEX 或 SET ... EX 需要的是從現在開始的持續時間
	durationUntilExpiry := time.Until(expiresAt)

	// 如果 token 已經過期或即將立即過期，則無需加入黑名單
	if durationUntilExpiry <= 0 {
		return nil // 或者可以記錄一下這個情況
	}

	// 構建 Redis 鍵
	key := tokenBlacklistPrefix + tokenString

	// 使用 SET 命令並設定 EX (秒為單位的過期時間)
	// Redis 會在 durationUntilExpiry 時間後自動刪除這個鍵
	// 值可以是任何簡單的標識，例如 "1" 或 true
	err := r.client.Set(ctx, key, "blacklisted", durationUntilExpiry).Err()
	if err != nil {
		// 實際應用中應記錄詳細錯誤
		// log.Printf("Error blacklisting token in Redis: %v", err)
		return err
	}

	return nil
}

// IsTokenBlacklisted 檢查指定的 token 字串是否存在於 Redis 黑名單中。
func (r *redisTokenBlacklistRepository) IsTokenBlacklisted(tokenString string) (bool, error) {
	ctx := context.Background() // 或者從呼叫者傳入 context

	key := tokenBlacklistPrefix + tokenString

	// 使用 EXISTS 命令檢查鍵是否存在
	// 如果鍵存在，表示 token 在黑名單中並且 Redis 尚未因 TTL 將其刪除
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		// 實際應用中應記錄詳細錯誤
		// log.Printf("Error checking token blacklist in Redis: %v", err)
		return false, err
	}

	// Exists 命令返回存在的鍵的數量，所以 > 0 表示存在
	return exists > 0, nil
}