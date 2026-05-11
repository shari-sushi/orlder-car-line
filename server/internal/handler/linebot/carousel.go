package linebot

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/line/line-bot-sdk-go/v8/linebot"

	"orlder-car-line/server/internal/infra/carsensor"
)

const carsensorBase = "https://www.carsensor.net"

type suvEntry struct {
	label      string // カルーセルのタイトル（日本語）
	imageText  string // スクレイピング失敗時のプレースホルダー文字列
	listingURL string // カーセンサー一覧ページ URL
}

var japanSUVs = []suvEntry{
	{"トヨタ ライズ", "RAIZE", carsensorBase + "/usedcar/bTO/s247/btX/index.html?SORT=3"},
	{"スバル XV", "SUBARU XV", carsensorBase + "/usedcar/bSB/s057/btX/index.html?SORT=3"},
	{"スズキ ジムニー", "JIMNY", carsensorBase + "/usedcar/bSZ/s010/btX/index.html?SORT=3"},
}

func placeholderImageURL(text string) string {
	// LINE は SVG 非対応のため .png を明示する。日本語は文字化けするので英語テキストを渡す。
	return fmt.Sprintf("https://placehold.co/1024x678/333333/ffffff.png?text=%s", url.QueryEscape(text))
}

func replyJapanSUVCarousel(replyToken string) error {
	type result struct {
		index int
		car   *carsensor.Car
	}

	results := make([]*carsensor.Car, len(japanSUVs))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, s := range japanSUVs {
		wg.Add(1)
		go func(idx int, entry suvEntry) {
			defer wg.Done()
			car, err := carsensor.FetchFirst(entry.listingURL)
			if err != nil {
				log.Printf("[linebot] %s 取得エラー: %v", entry.label, err)
				return
			}
			mu.Lock()
			results[idx] = car
			mu.Unlock()
		}(i, s)
	}
	wg.Wait()

	cols := make([]*linebot.CarouselColumn, 0, len(japanSUVs))
	for i, s := range japanSUVs {
		car := results[i]

		imageURL := placeholderImageURL(s.imageText)
		text := "情報取得中..."
		detailURL := s.listingURL

		if car != nil {
			if car.ImageURL != "" {
				imageURL = car.ImageURL
			}
			text = buildCarText(car)
			if car.DetailURL != "" {
				detailURL = car.DetailURL
			}
		}

		col := linebot.NewCarouselColumn(
			imageURL,
			s.label,
			text,
			linebot.NewURIAction("カーセンサーで見る", detailURL),
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

// buildCarText はグレード・価格・年式・走行距離をまとめる（LINE テキスト60文字制限に収める）。
func buildCarText(car *carsensor.Car) string {
	price := car.Price
	if price == "" {
		price = "価格不明"
	}
	year := car.Year
	if year == "" {
		year = "年式不明"
	}
	dist := car.Distance
	if dist == "" {
		dist = "距離不明"
	}
	base := fmt.Sprintf("%s / %s / %s", price, year, dist)

	// グレードは車名の先頭メーカー+車種を除いた部分（例: "ハイブリッド 1.2 G"）
	grade := extractGrade(car.Name)
	if grade == "" {
		return truncateRune(base, 60)
	}
	// grade + "\n" + base が 60文字以内に収まるよう grade を切り詰める
	maxGrade := 60 - 1 - len([]rune(base))
	if maxGrade <= 0 {
		return truncateRune(base, 60)
	}
	return fmt.Sprintf("%s\n%s", truncateRune(grade, maxGrade), base)
}

// truncateRune は s を最大 n ルーン文字に切り詰める。
func truncateRune(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// extractGrade は "トヨタ ライズ ハイブリッド 1.2 G 登録済未使用車..." のようなフルネームから
// メーカー・車種名（先頭2単語）を除いたグレード部分（最大3単語）を返す。
// 例: "ハイブリッド 1.2 G"
func extractGrade(name string) string {
	parts := strings.Fields(name)
	if len(parts) <= 2 {
		return ""
	}
	grade := parts[2:]
	if len(grade) > 3 {
		grade = grade[:3]
	}
	return strings.Join(grade, " ")
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

	imageURL := car.ImageURL
	if imageURL == "" {
		imageURL = placeholderImageURL("RAIZE")
	}

	msg := linebot.NewTemplateMessage(
		"トヨタ ライズ 最新情報",
		linebot.NewButtonsTemplate(
			imageURL,
			name,
			buildCarText(car),
			linebot.NewURIAction("カーセンサーで見る", car.DetailURL),
		),
	)
	_, err = getClient().ReplyMessage(replyToken, msg).Do()
	return err
}
