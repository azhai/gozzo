package mapper

import (
	"reflect"
	"sort"
)

// GetIndirectType 获取对象（指针）的实际类型
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

// GetFinalType 获取指针的实际类型，或数组哈希表的元素类型
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

// SortedMapKeys 对map的key按字母排序
func SortedMapKeys(data any) (keys []string) {
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
