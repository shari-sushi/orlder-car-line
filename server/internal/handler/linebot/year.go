package linebot

import "strings"

// YearRange は年式の範囲を表す。
type YearRange struct {
	Min int // 0 = 下限なし
	Max int // 0 = 上限なし
}

// ParseYear は年式文字列をパースして YearRange を返す。
//
//   - "2020" / "2020年"              → Max=2020
//   - "～2020" / "~2020"             → Max=2020
//   - "2018～" / "2018~"             → Min=2018
//   - "2018～2020" / "2018~2020"     → Min=2018, Max=2020
func ParseYear(s string) YearRange {
	s = normalizeYear(s)

	left, right, hasRange := splitRange(s)

	if !hasRange {
		return YearRange{Max: parseJapaneseInt(left)}
	}

	var min, max int
	if left != "" {
		min = parseJapaneseInt(left)
	}
	if right != "" {
		max = parseJapaneseInt(right)
	}
	return YearRange{Min: min, Max: max}
}

func normalizeYear(s string) string {
	s = strings.TrimSpace(s)
	for _, suffix := range []string{"年式", "年"} {
		s = strings.TrimSuffix(s, suffix)
	}
	return strings.TrimSpace(s)
}
