// config パッケージ: 環境変数の一元管理
// app/api/_lib/env.ts に相当
package config

import (
	"log"
	"os"
)

func getRequired(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("[警告] 環境変数 %s が未設定です", key)
	}
	return v
}

func getOptional(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

var (
	Port = getOptional("PORT", "8080")

	// LINE Messaging API
	LineChannelSecret      = getRequired("LINE_CHANNEL_SECRET")
	LineChannelAccessToken = getRequired("LINE_CHANNEL_ACCESS_TOKEN")

	// 認証 (例: "user1:pass1,user2:pass2")
	AuthUserPass = getRequired("AUTH_USER_PASS")

	// Redis: STORAGE_URL → REDIS_URL → デフォルト の優先順
	RedisURL = getOptional("STORAGE_URL", getOptional("REDIS_URL", "redis://localhost:6379"))
)
