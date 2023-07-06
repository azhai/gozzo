package match

import (
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// 字符串比较方式
const (
	CmpStringOmit            = iota // 不比较
	CmpStringContains               // 包含
	CmpStringStartswith             // 打头
	CmpStringEndswith               // 结尾
	CmpStringIgnoreSpaces           // 忽略空格
	CmpStringCaseInsensitive        // 不分大小写
	CmpStringEqual                  // 相等
)

var reDigit = regexp.MustCompile(`^\d+$`)

// StringMatch 比较是否相符
func StringMatch(a, b string, cmp int) bool {
	switch cmp {
	case CmpStringOmit:
		return true
	case CmpStringContains:
		return strings.Contains(a, b)
	case CmpStringStartswith:
		return strings.HasPrefix(a, b)
	case CmpStringEndswith:
		return strings.HasSuffix(a, b)
	case CmpStringIgnoreSpaces:
		a, b = RemoveSpaces(a), RemoveSpaces(b)
		return strings.EqualFold(a, b)
	case CmpStringCaseInsensitive:
		return strings.EqualFold(a, b)
	default: // 包括 CMP_STRING_EQUAL
		return strings.Compare(a, b) == 0
	}
}

// 是否在字符串列表中，只适用于CMP_STRING_EQUAL和CMP_STRING_STARTSWITH
func compareStringList(x string, lst []string, cmp int) bool {
	size := len(lst)
	if size == 0 {
		return false
	}
	if !sort.StringsAreSorted(lst) {
		sort.Strings(lst)
	}
	i := sort.Search(size, func(i int) bool { return lst[i] >= x })
	return i < size && StringMatch(x, lst[i], cmp)
}

// InStringList 是否在字符串列表中
func InStringList(x string, lst []string) bool {
	return compareStringList(x, lst, CmpStringEqual)
}

// StartStringList 是否在字符串列表中，比较方式是有任何一个开头符合
func StartStringList(x string, lst []string) bool {
	return compareStringList(x, lst, CmpStringStartswith)
}

// IsSubsetList lst1 是否 lst2 的（真）子集
func IsSubsetList(lst1, lst2 []string, strict bool) bool {
	if len(lst1) > len(lst2) {
		return false
	}
	if strict && len(lst1) == len(lst2) {
		return false
	}
	for _, x := range lst1 {
		if !InStringList(x, lst2) {
			return false
		}
	}
	return true
}

// IsDigit 是否纯数字
func IsDigit(s string) bool {
	return reDigit.MatchString(s)
}

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
