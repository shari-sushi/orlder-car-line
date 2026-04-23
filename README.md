# web + Discord bot 用テンプレート

React + Vercel Functions + Redis で動く Web アプリと Discord Bot を一体化したテンプレートです。

## 技術スタック

| 項目           | 内容                                        |
| -------------- | ------------------------------------------- |
| フロントエンド | React + TypeScript + Vite                   |
| スタイリング   | Tailwind CSS v4                             |
| ルーティング   | React Router v7                             |
| API            | Vercel Functions（`app/api/` ディレクトリ） |
| データベース   | Redis（Upstash 等）                         |
| Bot            | Discord Interactions API                    |

## セットアップ

### 前提条件

- Node.js 22 以上
- Docker（ローカル Redis 用）または Upstash 等の Redis サービス

### 1. リポジトリのクローンと依存関係のインストール

```bash
git clone <このリポジトリの URL>
cd <リポジトリ名>/app
npm install
```

### 2. Discord アプリの作成

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. **「New Application」** でアプリを作成
3. **「Bot」** タブ → **「Reset Token」** でトークンを取得・コピー
4. **「General Information」** タブ → **「Public Key」** をコピー
5. **「OAuth2」** タブ → **「CLIENT ID」**（Application ID）をコピー

### 3. 環境変数の設定

`app/.env.local` を作成し、以下を記載します（`.gitignore` 対象のため Git には含まれません）。

```sh
# Discord
DISCORD_PUBLIC_KEY=<Public Key>
DISCORD_APPLICATION_ID=<Application ID / CLIENT ID>
DISCORD_BOT_TOKEN=<Bot Token>

# 開発用ギルドコマンド（省略するとグローバルコマンドになり反映に最大1時間かかる）
DISCORD_COMMAND_GUILD_ID=<サーバーID>

# 認証（例: user1:pass1,user2:pass2）
AUTH_USER_PASS=admin:0000

# Redis 接続 URL（未設定時は redis://localhost:6379 にフォールバック）
REDIS_URL=redis://localhost:6379
```

### 4. ローカル Redis の起動

```bash
# プロジェクトルートで実行
docker-compose up -d redis
```

### 5. 開発サーバーの起動

```bash
cd app
npm run dev
# → http://localhost:5173/
```

## ローカルでの Discord Bot 動作確認（ngrok）

Discord の Interactions はサーバーへの HTTP リクエストで届くため、ローカルサーバーを ngrok で公開する必要があります。

```bash
# ngrok のインストール（未インストールの場合）
npm install -g ngrok

# ローカルサーバーを起動しつつ ngrok でトンネルを開く
make ngrok-up
```

起動後、ngrok が発行する URL（例: `https://xxxx.ngrok.io`）を Discord Developer Portal の
**「Interactions Endpoint URL」** に `https://xxxx.ngrok.io/api/discord` として設定します。

> ngrok の無料プランでは URL がセッションごとに変わります。開発中は毎回更新が必要です。

---

## Discord Interactions URL の設定

Vercel にデプロイ後、Discord Developer Portal で Interactions URL を設定します。

1. [Discord Developer Portal](https://discord.com/developers/applications) → アプリを選択
2. **「General Information」** タブ
3. **「Interactions Endpoint URL」** に `https://<your-app>.vercel.app/api/discord` を入力
4. **「Save Changes」**

## Vercel へのデプロイ

### プロジェクトのインポート

1. [Vercel](https://vercel.com/) にログイン → **「Add New → Project」**
2. このリポジトリを選択 → **「Import」**
3. 設定:
   - **Root Directory**: `app`
   - **Framework Preset**: `Vite`（自動検出）
4. **「Deploy」**

### 環境変数の登録

Vercel ダッシュボード → **「Settings」→「Environment Variables」** に以下を追加します。

| 変数名                     | 必須 | 説明                                                |
| -------------------------- | ---- | --------------------------------------------------- |
| `DISCORD_PUBLIC_KEY`       | ○    | Discord アプリの公開鍵                              |
| `DISCORD_APPLICATION_ID`   | ○    | Discord アプリケーション ID                         |
| `DISCORD_BOT_TOKEN`        | ○    | Discord Bot トークン                                |
| `DISCORD_COMMAND_GUILD_ID` | -    | コマンド登録先ギルド ID（省略時はグローバル）       |
| `AUTH_USER_PASS`           | ○    | 認証情報（例: `user1:pass1,user2:pass2`）           |
| `STORAGE_URL`              | -    | Redis 接続 URL（優先）。未設定時は `REDIS_URL` 参照 |
| `REDIS_URL`                | -    | Redis 接続 URL（例: Upstash の TLS URL）            |

### Discord スラッシュコマンドの登録

本番 URL が確定したら、以下を実行してコマンドを登録します。

```bash
cd app
VERCEL_ENV=production npx tsx api/discord/commands/register.ts
```

> `DISCORD_COMMAND_GUILD_ID` が設定されている場合はギルドコマンド（即時反映）、未設定の場合はグローバルコマンド（最大1時間）として登録されます。

## カスタマイズ

### スラッシュコマンドの追加

1. [app/api/discord/commands/](app/api/discord/commands/) にコマンド処理ファイルを追加
2. [app/api/discord/commands/register.ts](app/api/discord/commands/register.ts) の `commandDefs` にコマンド定義を追記
3. [app/api/discord.ts](app/api/discord.ts) でコマンド名に応じてハンドラを呼び出す処理を追加

### Web API の追加

`app/api/web/` 配下にファイルを追加します。Vercel Functions のルーティングは `app/vercel.json` で管理しています。

---

## ディレクトリ構成

```txt
.
├── app/
│   ├── api/                    # Vercel Functions（サーバーレス API）
│   │   ├── _lib/               # 認証・Redis・セッション等のユーティリティ
│   │   ├── discord/            # Discord Interactions エンドポイント
│   │   │   └── commands/       # スラッシュコマンド実装
│   │   └── web/                # Web フロントエンド向け API
│   ├── src/
│   │   ├── _domains/           # ドメインロジック
│   │   ├── components/         # 共通コンポーネント
│   │   ├── context/            # React Context
│   │   ├── hooks/              # カスタム hooks
│   │   ├── pages/              # ページコンポーネント
│   │   └── utils/              # ユーティリティ関数
│   └── vercel.json             # Vercel ルーティング設定
├── docker-compose.yml          # ローカル Redis 用
└── Makefile
```
