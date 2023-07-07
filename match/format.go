package match

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// RemoveSpaces 删除所有空白，包括中间的
func RemoveSpaces(s string) string {
	subs := map[string]string{
		" ": "", "\n": "", "\r": "", "\t": "", "\v": "", "\f": "",
	}
	return ReplaceWith(s, subs)
}

// ReduceSpaces 将多个连续空白缩减为一个空格
func ReduceSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ReplaceWith 一一对应进行替换，次序不定（因为map的关系）
func ReplaceWith(s string, subs map[string]string) string {
	if s == "" {
		return ""
	}
	var marks []string
	for key, value := range subs {
		marks = append(marks, key, value)
	}
	replacer := strings.NewReplacer(marks...)
	return replacer.Replace(s)
}

// TruncateText 截断长文本
func TruncateText(msg string, size int) string {
	if size <= 3 || len(msg) <= size {
		return msg
	}
	// 可能含有中文，要以rune计算结尾
	for i := size - 3; i >= 0; i-- {
		if utf8.RuneStart(msg[i]) {
			return msg[:i] + "..."
		}
	}
	return "..."
}

// QuoteColumns 盲转义，认为字段名以小写字母开头
func QuoteColumns(cols []string, sep string, quote func(string) string) string {
	re := regexp.MustCompile("([a-z][a-zA-Z0-9_]+)")
	repl, origin := quote("$1"), strings.Join(cols, sep)
	result := re.ReplaceAllString(origin, repl)
	if pad := (len(repl) - len("$1")) / 2; pad > 0 {
		left, right := repl[:pad], repl[len(repl)-pad:]
		oldNews := []string{
			left + left, left, right + right, right,
			"'" + left, "'", left + "'", "'",
		}
		result = strings.NewReplacer(oldNews...).Replace(result)
	}
	return result
}
