# Docker 構成計画

実装開始時に対応する。

## 方針

- `docker compose up` 一発で Redis + Go ボットが起動する構成にする
- ngrok はローカル開発時のみ使用

## docker-compose.yml（予定）

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  bot:
    build: ./bot
    ports:
      - "8080:8080"
    env_file: ./bot/.env.local
    depends_on:
      - redis
```

## bot/Dockerfile（予定）

Go のマルチステージビルドで軽量イメージを作る。

```dockerfile
# ビルドステージ
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server

# 実行ステージ
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

## VPS 移行時（さくら / ConoHa 等）

1. VPS に Docker / Docker Compose をインストール
2. リポジトリをクローン or バイナリを転送
3. `.env` を配置して `docker compose up -d`
4. LINE Developer Console の Webhook URL を更新

ポート開放と HTTPS 終端（nginx + Let's Encrypt 等）が必要になる点に注意。
