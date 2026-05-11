package linebot

import (
	"fmt"
	"log"
	"net/url"

	"github.com/line/line-bot-sdk-go/v8/linebot"

	"orlder-car-line/server/internal/infra/carsensor"
)

const carsensorBase = "https://www.carsensor.net"

type suvEntry struct {
	label      string // カルーセルのタイトル
	listingURL string // カーセンサー一覧ページ URL
}

var japanSUVs = []suvEntry{
	{"トヨタ ライズ", carsensorBase + "/usedcar/bTO/s247/btX/index.html?SORT=3"},
	{"スバル XV", carsensorBase + "/usedcar/bSB/s057/btX/index.html?SORT=3"},
	{"スズキ ジムニー", carsensorBase + "/usedcar/bSZ/s010/btX/index.html?SORT=3"},
}

func placeholderImageURL(label string) string {
	return fmt.Sprintf("https://placehold.co/1024x768/333333/ffffff?text=%s", url.QueryEscape(label))
}

func replyJapanSUVCarousel(replyToken string) error {
	cols := make([]*linebot.CarouselColumn, 0, len(japanSUVs))
	for _, s := range japanSUVs {
		col := linebot.NewCarouselColumn(
			placeholderImageURL(s.label),
			s.label,
			"カーセンサーで中古車を探す",
			linebot.NewURIAction("検索する", s.listingURL),
		)
		cols = append(cols, col)
	}

	msg := linebot.NewTemplateMessage(
		"国産SUV一覧",
		linebot.NewCarouselTemplate(cols...),
	)
	_, err := getClient().ReplyMessage(replyToken, msg).Do()
	return err
}

const raizeListingURL = carsensorBase + "/usedcar/bTO/s247/btX/index.html?SORT=3"

func replyLatestRaize(replyToken string) error {
	car, err := carsensor.FetchFirst(raizeListingURL)
	if err != nil {
		log.Printf("[linebot] ライズ取得エラー: %v", err)
		return reply(replyToken, "最新情報の取得に失敗しました")
	}

	name := car.Name
	if name == "" {
		name = "トヨタ ライズ"
	}

	price := car.Price
	if price == "" {
		price = "価格情報なし"
	}

	imageURL := car.ImageURL
	if imageURL == "" {
		imageURL = placeholderImageURL("ライズ")
	}

	msg := linebot.NewTemplateMessage(
		"トヨタ ライズ 最新情報",
		linebot.NewButtonsTemplate(
			imageURL,
			name,
			price,
			linebot.NewURIAction("カーセンサーで見る", car.DetailURL),
		),
	)
	_, err = getClient().ReplyMessage(replyToken, msg).Do()
	return err
}
