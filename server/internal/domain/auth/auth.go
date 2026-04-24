// auth パッケージ: Authorization ヘッダーの検証
// app/api/_lib/auth.ts に相当
// Bearer（セッショントークン）と Basic（curl 向け）の両方をサポート
package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"orlder-car-line/server/internal/infra/config"
	"orlder-car-line/server/internal/domain/session"
)

// Result は認証結果を表す。
type Result struct {
	Valid    bool
	Username string
	Err      string
}

// ParseUsers は AUTH_USER_PASS 環境変数をパースして map[username]password を返す。
func ParseUsers() map[string]string {
	users := map[string]string{}
	for _, entry := range strings.Split(config.AuthUserPass, ",") {
		idx := strings.Index(entry, ":")
		if idx == -1 {
			continue
		}
		user := strings.TrimSpace(entry[:idx])
		pass := strings.TrimSpace(entry[idx+1:])
		if user != "" {
			users[user] = pass
		}
	}
	return users
}

// ValidateHeader は Authorization ヘッダーを検証する。
func ValidateHeader(ctx context.Context, authHeader string) Result {
	if authHeader == "" {
		return Result{Err: "認証ヘッダーが必要です"}
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		data, err := session.Validate(ctx, token)
		if err != nil || data == nil {
			return Result{Err: "無効な認証トークンです"}
		}
		return Result{Valid: true, Username: data.Username}
	}

	if strings.HasPrefix(authHeader, "Basic ") {
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			return Result{Err: "不正な Basic 認証形式です"}
		}
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return Result{Err: "ユーザー名とパスワードが必要です"}
		}
		users := ParseUsers()
		if users[parts[0]] != parts[1] {
			return Result{Err: "無効なユーザー名またはパスワードです"}
		}
		return Result{Valid: true, Username: parts[0]}
	}

	return Result{Err: "Bearer または Basic を使用してください"}
}