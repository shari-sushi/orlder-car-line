// carsensor パッケージ: カーセンサー一覧ページから車両情報をスクレイピングする
package carsensor

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Car は一覧ページから取得した車両情報を表す。
type Car struct {
	Name     string
	Price    string
	ImageURL string
	DetailURL string
}

var (
	reDetailPath = regexp.MustCompile(`/usedcar/detail/(AU\d+)/index\.html`)
	reImage      = regexp.MustCompile(`ccsrpcma\.carsensor\.net/CSphoto/[^"'\s]+\.JPG`)
	rePrice      = regexp.MustCompile(`(\d+(?:\.\d+)?)万円`)
	reCarName    = regexp.MustCompile(`class="[^"]*CarName[^"]*"[^>]*>([^<]+)<`)
)

// FetchFirst はページURLから最初の1台を取得して返す。
func FetchFirst(listingURL string) (*Car, error) {
	body, err := fetchHTML(listingURL)
	if err != nil {
		return nil, err
	}

	car := &Car{}

	// 詳細ページパス
	if m := reDetailPath.FindString(body); m != "" {
		car.DetailURL = "https://www.carsensor.net" + m
	}

	// 画像URL（最初にマッチしたもの）
	if m := reImage.FindString(body); m != "" {
		car.ImageURL = "https://" + m
	}

	// 価格
	if m := rePrice.FindStringSubmatch(body); len(m) > 1 {
		car.Price = m[1] + "万円"
	}

	// 車名（CarName クラス要素 → なければ詳細URLから ID のみ）
	if m := reCarName.FindStringSubmatch(body); len(m) > 1 {
		car.Name = strings.TrimSpace(m[1])
	}

	if car.DetailURL == "" {
		return nil, fmt.Errorf("carsensor: 車両情報が取得できませんでした")
	}

	return car, nil
}

func fetchHTML(targetURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", err
	}
	// UA を設定しないと弾かれることがある
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; orlder-car-line-bot/1.0)")
	req.Header.Set("Accept-Language", "ja,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
