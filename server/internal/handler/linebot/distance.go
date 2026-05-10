package linebot

import (
	"fmt"
	"strconv"
	"strings"
)

// DistanceRange は走行距離の範囲を表す（単位: km）。
type DistanceRange struct {
	Min int // 0 = 下限なし
	Max int // 0 = 上限なし
}

func (d DistanceRange) String() string {
	switch {
	case d.Min > 0 && d.Max > 0:
		return fmt.Sprintf("%dkm以上%dkm以下", d.Min, d.Max)
	case d.Min > 0:
		return fmt.Sprintf("%dkm以上", d.Min)
	case d.Max > 0:
		return fmt.Sprintf("%dkm以下", d.Max)
	default:
		return ""
	}
}

// ParseDistance は距離文字列をパースして DistanceRange を返す。
//
//   - "1000" / "1千" / "1000km"      → Max=1000
//   - "～1000" / "~1000"             → Max=1000
//   - "1000～" / "1000~"             → Min=1000
//   - "1000～2000" / "1000~2000"     → Min=1000, Max=2000
func ParseDistance(s string) DistanceRange {
	s = normalizeDistance(s)

	left, right, hasRange := splitRange(s)

	if !hasRange {
		return DistanceRange{Max: parseJapaneseInt(s)}
	}

	var min, max int
	if left != "" {
		min = parseJapaneseInt(left)
	}
	if right != "" {
		max = parseJapaneseInt(right)
	}
	return DistanceRange{Min: min, Max: max}
}

// normalizeDistance は km 表記と前後の空白を除去する。
func normalizeDistance(s string) string {
	s = strings.TrimSpace(s)
	for _, unit := range []string{"km", "ｋｍ", "KM", "Km", "kM"} {
		s = strings.ReplaceAll(s, unit, "")
	}
	return strings.TrimSpace(s)
}

// parseJapaneseInt は "5万3千200" のような日本語混じり数字を int に変換する。
func parseJapaneseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	total := 0

	if before, after, found := strings.Cut(s, "万"); found {
		if n, err := strconv.Atoi(before); err == nil {
			total += n * 10000
		}
		s = after
	}

	if before, after, found := strings.Cut(s, "千"); found {
		if n, err := strconv.Atoi(before); err == nil {
			total += n * 1000
		}
		s = after
	}

	if s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			total += n
		}
	}

	return total
}
