// linebot パッケージ: LINE Messaging API Webhook ハンドラ
package linebot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"

	"orlder-car-line/server/internal/infra/config"
)

var (
	client     *linebot.Client
	clientOnce sync.Once
)

func getClient() *linebot.Client {
	clientOnce.Do(func() {
		var err error
		client, err = linebot.New(config.LineChannelSecret, config.LineChannelAccessToken)
		if err != nil {
			log.Fatalf("[linebot] クライアントの初期化に失敗: %v", err)
		}
	})
	return client
}

// HandleWebhook は LINE プラットフォームからの Webhook を受け取る。
// POST /webhook
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("[linebot] webhook received")

	cb, err := webhook.ParseRequest(config.LineChannelSecret, r)
	if err != nil {
		log.Printf("[linebot] ParseRequest エラー: %v", err)
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	ctx := r.Context()
	for _, event := range cb.Events {
		if e, ok := event.(webhook.MessageEvent); ok {
			if msg, ok := e.Message.(webhook.TextMessageContent); ok {
				log.Printf("[linebot] テキスト受信: %s", msg.Text)
				if err := handleMessage(ctx, e.ReplyToken, e.Source, msg.Text); err != nil {
					log.Printf("[linebot] handleMessage エラー: %v", err)
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func handleMessage(ctx context.Context, replyToken string, source webhook.SourceInterface, text string) error {
	userID := getUserID(source)
	if userID == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(text), "\n")

	// 3行以上: 車種・年式・走行距離を一括入力
	if len(lines) >= 3 {
		carName := strings.TrimSpace(lines[0])
		year := normalizeYear(strings.TrimSpace(lines[1]))
		distance := ParseDistance(strings.TrimSpace(lines[2]))
		url := buildSearchURL(carName, year, distance)
		return reply(replyToken, fmt.Sprintf("検索結果はこちら:\n%s", url))
	}

	// 1行: 使い方を案内
	return reply(replyToken, "フリーワード(車種やメーカー)\n年式\n距離\nを3行に分けて入力してください")
}

func normalizeYear(year string) string {
	return strings.TrimSuffix(year, "年")
}

func getUserID(source webhook.SourceInterface) string {
	switch s := source.(type) {
	case webhook.UserSource:
		return s.UserId
	case webhook.GroupSource:
		return s.UserId
	case webhook.RoomSource:
		return s.UserId
	}
	return ""
}

func reply(replyToken, text string) error {
	_, err := getClient().ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do()
	return err
}

func buildSearchURL(carName, year string, distance DistanceRange) string {
	params := url.Values{}
	params.Set("KW", carName)
	if year != "" {
		params.Set("YMIN", year)
		params.Set("YMAX", year)
	}
	if distance.Max > 0 {
		params.Set("SMAX", fmt.Sprintf("%d", distance.Max))
	}
	if distance.Min > 0 {
		params.Set("SMIN", fmt.Sprintf("%d", distance.Min))
	}
	return "https://www.carsensor.net/usedcar/search.php?" + params.Encode()
}
