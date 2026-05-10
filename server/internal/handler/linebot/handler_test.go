package linebot

import (
	"net/url"
	"regexp"
	"strings"
	"testing"
)

var reNumericOnly = regexp.MustCompile(`^\d+$`)

// TestBuildSearchURL はURLパラメータの形式を検証する。
func TestBuildSearchURL(t *testing.T) {
	tests := []struct {
		name     string
		carName  string
		year     string
		distance DistanceRange
		wantKW   string
		wantYMIN string
		wantYMAX string
		wantSMIN string
		wantSMAX string
	}{
		{
			name:     "車種・年式・上限距離あり",
			carName:  "xv",
			year:     "2020",
			distance: ParseDistance("5000"),
			wantKW:   "xv",
			wantYMIN: "2020",
			wantYMAX: "2020",
			wantSMAX: "5000",
		},
		{
			name:     "走行距離範囲指定",
			carName:  "フィット",
			year:     "2019",
			distance: ParseDistance("1万～5万"),
			wantKW:   "フィット",
			wantYMIN: "2019",
			wantYMAX: "2019",
			wantSMIN: "10000",
			wantSMAX: "50000",
		},
		{
			name:     "走行距離なし",
			carName:  "プリウス",
			year:     "2021",
			distance: DistanceRange{},
			wantKW:   "プリウス",
			wantYMIN: "2021",
			wantYMAX: "2021",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawURL := buildSearchURL(tt.carName, tt.year, tt.distance)

			parsed, err := url.Parse(rawURL)
			if err != nil {
				t.Fatalf("URLのパースに失敗: %v", err)
			}

			q := parsed.Query()

			// KW パラメータの確認
			if got := q.Get("KW"); got != tt.wantKW {
				t.Errorf("KW = %q, want %q", got, tt.wantKW)
			}

			// YMIN・YMAX・SMIN・SMAX は数字のみであること
			for _, key := range []string{"YMIN", "YMAX", "SMIN", "SMAX"} {
				if v := q.Get(key); v != "" && !reNumericOnly.MatchString(v) {
					t.Errorf("%s = %q: 数字以外の文字が含まれている", key, v)
				}
			}

			// 各パラメータの値確認
			check := func(key, want string) {
				t.Helper()
				got := q.Get(key)
				if want == "" && got != "" {
					t.Errorf("%s = %q, want 空（パラメータ不要）", key, got)
				} else if want != "" && got != want {
					t.Errorf("%s = %q, want %q", key, got, want)
				}
			}
			check("YMIN", tt.wantYMIN)
			check("YMAX", tt.wantYMAX)
			check("SMIN", tt.wantSMIN)
			check("SMAX", tt.wantSMAX)

			// URLに日本語の単位が含まれていないこと
			if strings.ContainsAny(rawURL, "以下以上万千") {
				t.Errorf("URLに日本語の単位が含まれている: %s", rawURL)
			}
		})
	}
}
