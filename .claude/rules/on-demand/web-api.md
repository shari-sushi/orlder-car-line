---
paths: ["適応させたいファイルpath"]
---

# Web API・認証ルール（テンプレート）

## Web API認証（`/api/web/*` エンドポイント）

## 方式: Bearer Token認証（ブラウザ向け）

ログインしてセッショントークンを取得し、以降のリクエストで使用します。

**認証フロー:**

1. `/api/web/auth/login` にユーザー名とパスワードを送信
2. セッショントークンを取得
3. 以降のリクエストで `Authorization: Bearer {token}` ヘッダーに付与

**セッション仕様:**

- **有効期限**: 7日間
- **自動延長**: API使用のたびに延長（トークンは不変）
- **保存場所**: Redis（キー: `session:{token}`）

**使用例（curl）:**

```bash
# ログイン
curl -X POST /api/web/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'

# 認証付きリクエスト
curl /api/web/protected \
  -H "Authorization: Bearer {token}"
```
