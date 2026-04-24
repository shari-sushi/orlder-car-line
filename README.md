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
1. ユーザーが車種名をテキストで打ち込む
2. 類似名検索して選択肢をクイックリプライで返す
3. ユーザーが選択
4. 年式を聞く
5. 走行距離を聞く
6. カーセンサーの検索結果 URL を Flex Message で返す
```

## 技術スタック

| 項目           | 選定          | 備考                                              |
| -------------- | ------------- | ------------------------------------------------- |
| 言語           | Go            | 確定                                              |
| フレームワーク | 未定          | 入社先に確認中 → [比較メモ](docs/go-framework.md) |
| セッション管理 | Redis         | TTL 付き → [設計メモ](docs/redis-session.md)      |
| LINE           | Messaging API | → [概要メモ](docs/line-messaging-api.md)          |

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
- [ ] フレームワーク確認（入社先へ確認中）
- [ ] LINE Developers アカウント作成
- [ ] LINE 公式アカウント作成・Messaging API 有効化
- [ ] Channel Secret / Channel Access Token 取得
- [ ] Go ローカルサーバー立てる
- [ ] ngrok で Webhook 疎通確認
- [ ] Redis 接続実装
- [ ] 会話フロー実装
