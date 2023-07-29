package datatype

import (
	"reflect"
	"strconv"
)

// Number 数字类型
type Number interface {
	Integer | Unsigned | Float
}

// Integer 整数
type Integer interface {
	int | int64 | int32 | int16 | int8
}

// FormatInteger 整数转为字符串
func FormatInteger[T Integer](n T) string {
	return strconv.FormatInt(int64(n), 10)
}

// ParseInteger 字符串转为整数，可能失败
func ParseInteger[T Integer](s string) (T, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	return T(n), err
}

// GetInteger 读取整数值
func GetInteger[T Integer](value reflect.Value) T {
	return T(value.Int())
}

// Unsigned 无符号数
type Unsigned interface {
	uint | uint64 | uint32 | uint16 | uint8
}

// FormatUnsigned 无符号数转为字符串
func FormatUnsigned[T Unsigned](u T) string {
	return strconv.FormatUint(uint64(u), 10)
}

// ParseUnsigned 字符串转为无符号数，可能失败
func ParseUnsigned[T Unsigned](s string) (T, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	return T(u), err
}

// GetUnsigned 读取无符号数值
func GetUnsigned[T Unsigned](value reflect.Value) T {
	return T(value.Uint())
}

// Float 浮点数
type Float interface {
	float64 | float32
}

// FormatFloat 浮点数转为字符串
func FormatFloat[T Float](f T) string {
	return strconv.FormatFloat(float64(f), 'G', -1, 64)
}

// ParseFloat 字符串转为浮点数，可能失败
func ParseFloat[T Float](s string) (T, error) {
	f, err := strconv.ParseFloat(s, 64)
	return T(f), err
}

// GetFloat 读取浮点数值
func GetFloat[T Float](value reflect.Value) T {
	return T(value.Float())
}
