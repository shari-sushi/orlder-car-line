# LINE チャネルの種類

LINE Developers コンソールで作成できるチャネルは主に 4 種類ある。

## 1. Messaging API

Bot としてユーザーとメッセージのやり取りを行うためのチャネル。

- ユーザーからのメッセージを Webhook で受信
- テキスト・画像・Flex Message などを返信・プッシュ送信
- Channel Secret（署名検証）と Channel Access Token（API 認証）が必要

**主な用途**: 注文受付 Bot、通知 Bot、カスタマーサポート

## 2. LINE Login

ウェブサービスやアプリへ LINE アカウントで OAuth 2.0 ログインするためのチャネル。

- OpenID Connect に対応（プロフィール・メールアドレスの取得が可能）
- ログイン後に取得した `userId` を Messaging API の Push 送信先として利用できる
- LIFF と組み合わせることが多い

**主な用途**: 会員登録の簡略化、LINE アカウントとサービスアカウントの紐付け

## 3. LIFF（LINE Front-end Framework）

LINE アプリ内のブラウザ（WebView）で動作する Web アプリを作るためのチャネル。

- LINE Login チャネルに紐付けて作成する
- `liff.init()` でユーザー情報の取得やメッセージ送信 API を呼び出せる
- QR コードや URL で LINE 内から直接開ける

**主な用途**: 注文フォーム、アンケート、予約画面など LINE 内完結の Web UI

## 4. LINE Notify（廃止済み）

サーバーから LINE にプッシュ通知を送るための簡易チャネル。

- 2025 年 3 月末にサービス終了
- 代替として **Messaging API の Push Message** を使用する

**現在の対応**: Messaging API への移行が必要

---

## チャネル選択の目安

| やりたいこと                       | 使うチャネル  |
| ---------------------------------- | ------------- |
| Bot でメッセージを受け取りたい     | Messaging API |
| LINE アカウントでログインさせたい  | LINE Login    |
| LINE 内で Web フォームを表示したい | LIFF          |
| サーバーから通知を送りたい         | Messaging API |

## 参考

- [LINE Developers チャネル概要](https://developers.line.biz/ja/docs/liff/overview/)
- [Messaging API ドキュメント](https://developers.line.biz/ja/docs/messaging-api/)
- [LINE Login ドキュメント](https://developers.line.biz/ja/docs/line-login/)
- [LIFF ドキュメント](https://developers.line.biz/ja/docs/liff/)
