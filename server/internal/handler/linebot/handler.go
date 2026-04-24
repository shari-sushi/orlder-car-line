// linebot パッケージ: LINE Messaging API Webhook ハンドラ
// TODO: github.com/line/line-bot-sdk-go/v8 を使った実装
package linebot

import (
	"log"
	"net/http"
)

// HandleWebhook は LINE プラットフォームからの Webhook を受け取る。
// POST /webhook
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: 署名検証（X-Line-Signature ヘッダー）
	// TODO: イベント解析（メッセージ / フォローなど）
	// TODO: 会話フロー（Redis セッションを使ったステートマシン）
	log.Println("[linebot] webhook received")
	w.WriteHeader(http.StatusOK)
}
