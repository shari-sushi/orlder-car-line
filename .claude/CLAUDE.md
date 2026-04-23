# プロジェクト概要

人間が見えるすべての履歴やファイル作成は日本語で記述する。
詳細なルールは `.claude/rules/` 配下のファイルを参照する。

## 常時読み込み

@.claude/rules/coding-conventions.md

---

## Web アプリ

### 場所

`/app/`

### 起動

```bash
cd app
npm install
npm run dev
# → http://localhost:5173/
```

### 技術スタック

| 項目           | 内容                              |
| -------------- | --------------------------------- |
| フレームワーク | React + TypeScript                |
| ビルドツール   | Vite                              |
| スタイリング   | Tailwind CSS v4                   |
| ルーティング   | React Router                      |
| バックエンド   | Vercel Functions (`api/`) + Redis |

### ファイル構成

```txt
app/
├── api/                          # Vercel Functions（サーバーレス API）
│   ├── _lib/
│   │   ├── auth.ts               # 認証ユーティリティ
│   │   ├── discord/              # Discord API ラッパー
│   │   ├── env.ts                # 環境変数の一元管理
│   │   ├── redis.ts              # Redis クライアント singleton
│   │   ├── session.ts            # セッション管理
│   │   └── vercel-adapter.ts
│   ├── discord/                  # Discord Interactions エンドポイント
│   │   ├── _util/commands.ts
│   │   └── commands/             # スラッシュコマンド実装
│   ├── web/                      # Web フロントエンド向け API
│   │   ├── _handlers/
│   │   └── auth/
│   ├── discord.ts                # Discord Interactions エントリーポイント
│   ├── login.ts
│   ├── logout.ts
│   └── me.ts
├── src/
│   ├── _domains/                 # ドメインロジック
│   ├── components/               # 共通コンポーネント
│   ├── context/                  # React Context（auth/, menu/, overlay/）
│   ├── hooks/                    # カスタム hooks
│   ├── icons/                    # アイコンコンポーネント
│   ├── lib/                      # API クライアントなどライブラリ
│   ├── pages/                    # ページコンポーネント
│   ├── utils/                    # ユーティリティ関数
│   ├── App.tsx                   # BrowserRouter + Routes
│   └── main.tsx
└── ...設定ファイル
```

---

## バックエンド（Vercel Functions + Redis）

### 環境変数

| 変数                       | 必須 | 説明                                                  |
| -------------------------- | ---- | ----------------------------------------------------- |
| `DISCORD_PUBLIC_KEY`       | ○    | Discord アプリの公開鍵                                |
| `DISCORD_APPLICATION_ID`   | ○    | Discord アプリケーション ID                           |
| `DISCORD_BOT_TOKEN`        | ○    | Discord Bot トークン                                  |
| `DISCORD_COMMAND_GUILD_ID` | -    | コマンド登録先ギルド ID（省略時はグローバル）         |
| `AUTH_USER_PASS`           | ○    | 認証情報（例: `user1:pass1,user2:pass2`）             |
| `STORAGE_URL`              | -    | Redis 接続 URL（優先）。未設定時は `REDIS_URL` を参照 |
| `REDIS_URL`                | -    | Redis 接続 URL（例: `rediss://...` Upstash など）     |

ローカル開発時は `app/.env.develop` に記載。
Vercel 本番は Environment Variables に設定。
（`STORAGE_URL` / `REDIS_URL` ともに未設定の場合は `redis://localhost:6379` にフォールバック）

### 注意点

- 作業後、コンパイルエラーを確認し、発生していれば修正する
