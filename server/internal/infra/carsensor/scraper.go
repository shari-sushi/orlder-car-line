// carsensor パッケージ: カーセンサー一覧ページから車両情報をスクレイピングする
package carsensor

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// httpClient はカーセンサーへのリクエストに使うタイムアウト付きクライアント。
// LINE Webhook の応答期限（5秒）を考慮して 4秒 に設定する。
var httpClient = &http.Client{Timeout: 4 * time.Second}

// Car は一覧ページから取得した車両情報を表す。
type Car struct {
	Name      string
	Price     string
	Distance  string
	Year      string
	ImageURL  string
	DetailURL string
}

var (
	// 詳細ページURL: /usedcar/detail/AU1234567890/index.html
	reDetailPath = regexp.MustCompile(`/usedcar/detail/(AU\d+)/index\.html`)

	// 画像URL: document.write() の data-original 属性に入っている
	reImage = regexp.MustCompile(`data-original="(//ccsrpcma\.carsensor\.net/CSphoto/[^"\s>]+\.JPG)"`)

	// 車名: img の alt 属性（メーカー名で始まるもの）
	reAlt = regexp.MustCompile(`alt="((?:トヨタ|ホンダ|マツダ|日産|スバル|スズキ|三菱|ダイハツ|いすゞ)[^"]{1,300})"`)

	// 価格: basePrice__mainPriceNum スパンに分割されている整数部と小数部
	rePrice = regexp.MustCompile(`basePrice__mainPriceNum">(\d+)</span><span[^>]*>(\.\d+)</span>`)

	// 走行距離: specList の emphasisData スパン（km / 万km 両形式に対応）
	reDistance = regexp.MustCompile(`走行距離</dt>\s*<dd[^>]*><span[^>]*>([\d,\.]+)</span>(万km|km)`)

	// 年式: specList の emphasisData スパン
	reYear = regexp.MustCompile(`年式</dt>\s*<dd[^>]*><span[^>]*>(\d{4})</span>`)
)

// FetchFirst は一覧ページURLから最初の1台を取得して返す。
// ページ全体ではなく最初の detail URL 以降のブロックで絞り込む。
func FetchFirst(listingURL string) (*Car, error) {
	body, err := fetchHTML(listingURL)
	if err != nil {
		return nil, err
	}

	car := &Car{}

	// 最初の詳細ページURLを見つける
	loc := reDetailPath.FindStringIndex(body)
	if loc == nil {
		return nil, fmt.Errorf("carsensor: 車両情報が取得できませんでした")
	}
	car.DetailURL = "https://www.carsensor.net" + body[loc[0]:loc[1]]

	// 価格・年式・距離は detail URL より約7,000〜9,000文字後に出現するため
	// 前方 4000 文字・後方 10000 文字をブロックとして使う
	start := loc[0] - 4000
	if start < 0 {
		start = 0
	}
	end := loc[1] + 10000
	if end > len(body) {
		end = len(body)
	}
	block := body[start:end]

	// 画像URL
	if m := reImage.FindStringSubmatch(block); len(m) > 1 {
		car.ImageURL = "https:" + m[1]
	}

	// 車名（&nbsp; 等の HTML エンティティを除去して整形）
	if m := reAlt.FindStringSubmatch(block); len(m) > 1 {
		name := strings.ReplaceAll(m[1], "&nbsp;", " ")
		name = strings.ReplaceAll(name, "&amp;", "&")
		car.Name = strings.TrimSpace(name)
	}

	// 価格（整数部と小数部が別スパンに分割されているため結合）
	if m := rePrice.FindStringSubmatch(block); len(m) > 2 {
		car.Price = m[1] + m[2] + "万円"
	}

	// 走行距離（km / 万km 両形式）
	if m := reDistance.FindStringSubmatch(block); len(m) > 2 {
		car.Distance = m[1] + m[2]
	}

	// 年式
	if m := reYear.FindStringSubmatch(block); len(m) > 1 {
		car.Year = m[1] + "年"
	}

	return car, nil
}

func fetchHTML(targetURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "ja,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := httpClient.Do(req)
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
