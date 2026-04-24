# LINE Messaging API

## 基本構造

```txt
ユーザー
  ↕
LINE プラットフォーム
  ↕ Webhook (HTTPS/JSON)
自前の Go サーバー
```

LINE のWebhookはリクエストが来るたびに単発で叩かれる。
HTTP はステートレスなので前のやり取りをサーバーは覚えていない。

## 送れるメッセージ種類

| 種類                   | 用途                                 |
| ---------------------- | ------------------------------------ |
| テキスト               | 基本の返答                           |
| クイックリプライ       | 選択肢ボタン（メッセージ下部に表示） |
| Flex Message           | カード型 UI、自由レイアウト          |
| テンプレートメッセージ | カルーセル等の定型 UI                |

## メッセージの送り方

| 種類                | 説明                                               |
| ------------------- | -------------------------------------------------- |
| Reply               | ユーザーのメッセージへの返信（無料枠内）           |
| Push                | サーバーから任意タイミングで送信（査定結果通知等） |
| Broadcast/Multicast | 全員 or 複数人に一斉送信                           |

## 必要なもの

- **Channel Secret**: Webhook の署名検証に使う
- **Channel Access Token**: API リクエストの認証に使う
- **User ID**: Push Message の宛先

## Webhook 署名検証

LINE は Webhook リクエストに署名（`X-Line-Signature` ヘッダ）を付与する。
`Channel Secret` を使って HMAC-SHA256 で検証する。

```go
mac := hmac.New(sha256.New, []byte(channelSecret))
mac.Write(body)
signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
// signature == X-Line-Signature ヘッダの値であれば OK
```

## 参考

- [LINE Developers](https://developers.line.biz/ja/)
- [Messaging API リファレンス](https://developers.line.biz/ja/reference/messaging-api/)
