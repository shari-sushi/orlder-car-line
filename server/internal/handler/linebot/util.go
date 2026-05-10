package linebot

import "strings"

// splitRange は "A～B" / "A~B" 形式の文字列を左右に分割する。
// 区切りがない場合は left=s, right="", hasRange=false を返す。
func splitRange(s string) (left, right string, hasRange bool) {
	s = strings.ReplaceAll(s, "～", "~")
	return strings.Cut(s, "~")
}
