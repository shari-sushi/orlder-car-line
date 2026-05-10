# LINE中古車検索ボット

入社に向けて **LINE Messaging API** と **中古車売買ドメイン** の理解を深めるためのキャッチアップ用プロジェクト。

実用性よりも、技術・ドメインの習得にフォーカスする。

## 目的

参考サービス（LINE完結型中古車売買）の仕組みを自作することで、以下を習得する。

- LINE Messaging API（Webhook・クイックリプライ・Flex Message）
- ステートマシンによる会話フロー管理
- Redis を使ったセッション管理
- 中古車ドメインの基礎知識

## 作るもの

車種・年式・走行距離を LINE で入力すると、カーセンサーの検索結果 URL を Flex Message で返すボット。

### ユーザーフロー

```txt
1. ユーザーが車種・年式・走行距離を改行区切りで送る
2. Bot がカーセンサーをスクレイピングして価格帯を取得
3. 検索URL + 価格帯をユーザーに返信
4. ユーザーが価格帯を返信（例: 100万円〜200万円）
5. 価格帯で絞り込んだ URL を返す
```

## 技術スタック

| 項目           | 選定          | 備考                                         |
| -------------- | ------------- | -------------------------------------------- |
| 言語           | Go            | 確定                                         |
| フレームワーク | Gin           | → [比較メモ](docs/go-framework.md)           |
| セッション管理 | Redis         | TTL 付き → [設計メモ](docs/redis-session.md) |
| LINE           | Messaging API | → [概要メモ](docs/line-messaging-api.md)     |

## ディレクトリ構成

```txt
.
├── bot/                    # Go LINE ボット（実装予定）
│   ├── cmd/server/         # エントリーポイント
│   └── internal/           # ハンドラ・セッション・ドメインロジック
├── app/                    # Web 管理 UI（React + Vercel Functions）※テンプレート流用
├── docs/                   # 技術調査・選定メモ
├── docker-compose.yml      # ローカル Redis
└── Makefile
```

## ドキュメント

| ファイル                                                 | 内容                                          |
| -------------------------------------------------------- | --------------------------------------------- |
| [docs/line-messaging-api.md](docs/line-messaging-api.md) | LINE Messaging API の基本構造・メッセージ種類 |
| [docs/go-framework.md](docs/go-framework.md)             | Go フレームワーク比較・選定                   |
| [docs/redis-session.md](docs/redis-session.md)           | Redis セッション設計                          |

## 進捗

- [x] 仕様決め
- [x] 技術スタック選定（フレームワーク以外）
- [x] フレームワーク確認（Gin でいいや）
- [x] LINE Developers アカウント作成
- [x] LINE 公式アカウント作成・Messaging API 有効化
- [x] Channel Secret / Channel Access Token 取得
- [x] Go ローカルサーバー立てる
- [x] ngrok で Webhook 疎通確認
- [ ] Redis 接続実装
- [ ] カーセンサー価格帯スクレイピング
- [ ] ステートマシン実装（waiting_price）
- [ ] 価格帯パース・絞り込みURL生成
