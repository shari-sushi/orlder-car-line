// session パッケージ: Redis を使ったセッション管理
// app/api/_lib/session.ts に相当
package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	redisclient "orlder-car-line/server/internal/infra/redis"
)

const (
	prefix = "session:"
	TTL    = 36 * time.Hour
)

// Data はセッションに保存するデータ。
type Data struct {
	Username       string `json:"username"`
	CreatedAt      int64  `json:"createdAt"`
	LastAccessedAt int64  `json:"lastAccessedAt"`
}

// GenerateToken はランダムな 64 文字の hex トークンを生成する。
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Create はセッションを作成して Redis に保存し、トークンを返す。
func Create(ctx context.Context, username string) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}
	now := time.Now().UnixMilli()
	data := Data{
		Username:       username,
		CreatedAt:      now,
		LastAccessedAt: now,
	}
	if err := redisclient.Set(ctx, prefix+token, data, TTL); err != nil {
		return "", err
	}
	return token, nil
}

// Validate はトークンを検証し、有効なら最終アクセス時刻を更新して Data を返す。
func Validate(ctx context.Context, token string) (*Data, error) {
	if token == "" {
		return nil, nil
	}
	data, err := redisclient.Get[Data](ctx, prefix+token)
	if err != nil || data == nil {
		return nil, err
	}
	data.LastAccessedAt = time.Now().UnixMilli()
	if err := redisclient.Set(ctx, prefix+token, data, TTL); err != nil {
		return nil, err
	}
	return data, nil
}

// Delete はセッションを Redis から削除する。
func Delete(ctx context.Context, token string) error {
	return redisclient.Del(ctx, prefix+token)
}