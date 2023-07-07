package mapper

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

func GetIndirectType(v any) (rt reflect.Type) {
	var ok bool
	if rt, ok = v.(reflect.Type); !ok {
		rt = reflect.TypeOf(v)
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return
}

func GetFinalType(v any) (rt reflect.Type) {
	rt = GetIndirectType(v)
	for {
		switch rt.Kind() {
		default:
			return rt
		case reflect.Ptr, reflect.Chan:
			rt = rt.Elem()
		case reflect.Array, reflect.Slice:
			rt = rt.Elem()
		case reflect.Map:
			kk := rt.Key().Kind()
			if kk == reflect.String || kk <= reflect.Float64 {
				rt = rt.Elem()
			} else {
				return rt
			}
		}
	}
}

func SortedKeys(data any) (keys []string) {
	rt := GetIndirectType(data)
	if rt.Kind() != reflect.Map || rt.Key().Kind() != reflect.String {
		return // data 必须是 map[string]xxxx 类型
	}
	for _, rv := range reflect.ValueOf(data).MapKeys() {
		keys = append(keys, rv.String())
	}
	sort.Strings(keys)
	return
}

func GetColumns(v any, alias string, cols []string) []string {
	rt := GetIndirectType(v)
	if rt.Kind() != reflect.Struct {
		return cols
	}
	for i := 0; i < rt.NumField(); i++ {
		t := rt.Field(i).Tag.Get("json")
		if t == "" || t == "-" {
			continue
		} else if strings.HasSuffix(t, "inline") {
			cols = GetColumns(rt.Field(i).Type, alias, cols)
		} else {
			if alias != "" {
				t = fmt.Sprintf("%s.%s", alias, t)
			}
			cols = append(cols, t)
		}
	}
	return cols
}

func GetChangesFor(v any, changes map[string]any) map[string]any {
	result := make(map[string]any)
	cols := GetColumns(v, "", []string{})
	for _, c := range cols {
		if val, ok := changes[c]; ok {
			result[c] = val
		}
	}
	return result
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
