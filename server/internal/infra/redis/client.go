// redis パッケージ: Redis クライアント singleton + 汎用 Get/Set/Del
// app/api/_lib/redis.ts に相当
package redis

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"orlder-car-line/server/internal/infra/config"
)

var (
	once   sync.Once
	client *redis.Client
)

func getClient() *redis.Client {
	once.Do(func() {
		opt, err := redis.ParseURL(config.RedisURL)
		if err != nil {
			log.Fatalf("[Redis] URL のパースに失敗: %v", err)
		}
		client = redis.NewClient(opt)
	})
	return client
}

// Get は JSON デシリアライズして返す。キーが存在しない場合は nil, nil を返す。
func Get[T any](ctx context.Context, key string) (*T, error) {
	val, err := getClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Set は値を JSON シリアライズして保存する。ttl=0 は期限なし。
func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return getClient().Set(ctx, key, b, ttl).Err()
}

// Del はキーを削除する。
func Del(ctx context.Context, key string) error {
	return getClient().Del(ctx, key).Err()
}