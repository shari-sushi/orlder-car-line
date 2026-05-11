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

	// X-Line-Signature ヘッダーの値と、ChannelSecret で計算した署名が一致するか確認して安全性を担保(共有秘密鍵)
	// HMAC-SHA256(ChannelSecret, リクエストボディ) → 署名　により、bodyが改ざんされている場合も検出する
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
	// イベントドキュメント https://developers.line.biz/ja/reference/messaging-api/#webhook-event-objects
	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch msg := e.Message.(type) {
			case webhook.TextMessageContent:
				log.Printf("[linebot] テキスト受信: %s", msg.Text)
				if err := handleMessage(ctx, e.ReplyToken, e.Source, msg.Text); err != nil {
					log.Printf("[linebot] handleMessage エラー: %v", err)
				}
			case webhook.ImageMessageContent:
			case webhook.StickerMessageContent:
			}
		case webhook.FollowEvent:
		case webhook.UnfollowEvent:
		case webhook.PostbackEvent:
		case webhook.JoinEvent:
		case webhook.LeaveEvent:
		case webhook.MemberJoinedEvent:
		case webhook.MemberLeftEvent:
		}
	}

	w.WriteHeader(http.StatusOK)
}

func handleMessage(ctx context.Context, replyToken string, source webhook.SourceInterface, text string) error {
	userID := getUserID(source)
	if userID == "" {
		return nil
	}

	lower := strings.ToLower(strings.TrimSpace(text))

	// "suv" または "国産suv" → カルーセル表示
	if lower == "suv" || lower == "国産suv" || lower == "suvを探す" {
		return replyJapanSUVCarousel(replyToken)
	}

	// "ライズ" → カーセンサー最新1台を表示
	if lower == "ライズ" {
		return replyLatestRaize(replyToken)
	}

	// リッチメニュー: スタッフに相談する
	if lower == "スタッフに相談する" {
		return reply(replyToken, "ご相談はこちらまでお気軽にどうぞ！\n（スタッフ対応は準備中です）")
	}

	// リッチメニュー: 自分で探す → 検索の使い方を案内
	if lower == "自分で探す" {
		return reply(replyToken, "以下の3行を入力してください:\n\n例)\nハリアー\n2020\n5万")
	}

	lines := strings.Split(strings.TrimSpace(text), "\n")

	// 3行以上: 車種・年式・走行距離を一括入力
	if len(lines) >= 3 {
		carName := strings.TrimSpace(lines[0])
		year := ParseYear(strings.TrimSpace(lines[1]))
		distance := ParseDistance(strings.TrimSpace(lines[2]))
		url := buildSearchURL(carName, year, distance)
		return reply(replyToken, fmt.Sprintf("検索結果はこちら:\n%s", url))
	}

	// 1行: 使い方を案内
	return reply(replyToken, "フリーワード(車種やメーカー)\n年式\n距離\nを3行に分けて入力してください")
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

func buildSearchURL(carName string, year YearRange, distance DistanceRange) string {
	params := url.Values{}
	params.Set("KW", carName)
	if year.Min > 0 {
		params.Set("YMIN", fmt.Sprintf("%d", year.Min))
	}
	if year.Max > 0 {
		params.Set("YMAX", fmt.Sprintf("%d", year.Max))
	}
	if distance.Min > 0 {
		params.Set("SMIN", fmt.Sprintf("%d", distance.Min))
	}
	if distance.Max > 0 {
		params.Set("SMAX", fmt.Sprintf("%d", distance.Max))
	}
	return "https://www.carsensor.net/usedcar/search.php?" + params.Encode()
}
