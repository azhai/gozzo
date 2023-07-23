package mapper

import (
	"strconv"
	"time"
)

// Unit 数据转换单元
type Unit struct {
	Data []byte
}

// String 将字节码转为字符串
func (u *Unit) String() string {
	return string(u.Data)
}

// FromStr 将字符串转为字节码
func (u *Unit) FromStr(s string) error {
	u.Data = []byte(s)
	return nil
}

// ToStr 将字节码转为字符串
func (u *Unit) ToStr() (string, error) {
	return u.String(), nil
}

// FromInt 将整数转为字节码
func (u *Unit) FromInt(n int) error {
	return u.FromStr(strconv.Itoa(n))
}

// ToInt 将字节码转为整数
func (u *Unit) ToInt() (int, error) {
	return strconv.Atoi(u.String())
}

// FromInt64 将整数转为字节码
func (u *Unit) FromInt64(n int64) error {
	return u.FromStr(strconv.FormatInt(n, 10))
}

// ToInt64 将字节码转为整数
func (u *Unit) ToInt64() (int64, error) {
	return strconv.ParseInt(u.String(), 10, 64)
}

// FromInt32 将整数转为字节码
func (u *Unit) FromInt32(n int32) error {
	return u.FromInt64(int64(n))
}

// ToInt32 将字节码转为整数
func (u *Unit) ToInt32() (n int32, err error) {
	var x int64
	if x, err = u.ToInt64(); err == nil {
		n = int32(x)
	}
	return
}

// FromInt16 将整数转为字节码
func (u *Unit) FromInt16(n int16) error {
	return u.FromInt64(int64(n))
}

// ToInt16 将字节码转为整数
func (u *Unit) ToInt16() (n int16, err error) {
	var x int64
	if x, err = u.ToInt64(); err == nil {
		n = int16(x)
	}
	return
}

// FromInt8 将整数转为字节码
func (u *Unit) FromInt8(n int8) error {
	return u.FromInt64(int64(n))
}

// ToInt8 将字节码转为整数
func (u *Unit) ToInt8() (n int8, err error) {
	var x int64
	if x, err = u.ToInt64(); err == nil {
		n = int8(x)
	}
	return
}

// FromUint 将整数转为字节码
func (u *Unit) FromUint(n uint) error {
	return u.FromStr(strconv.FormatUint(uint64(n), 10))
}

// ToUint 将字节码转为整数
func (u *Unit) ToUint() (n uint, err error) {
	var x uint64
	if x, err = u.ToUint64(); err == nil {
		n = uint(x)
	}
	return
}

// FromUint64 将整数转为字节码
func (u *Unit) FromUint64(n uint64) error {
	return u.FromStr(strconv.FormatUint(n, 10))
}

// ToUint64 将字节码转为整数
func (u *Unit) ToUint64() (uint64, error) {
	return strconv.ParseUint(u.String(), 10, 64)
}

// FromFloat64 将浮点数转为字节码
func (u *Unit) FromFloat64(f float64) error {
	return u.FromStr(strconv.FormatFloat(f, 'G', -1, 64))
}

// ToFloat64 将字节码转为浮点数
func (u *Unit) ToFloat64() (f float64, err error) {
	return strconv.ParseFloat(u.String(), 64)
}

// FromBool 将布尔值转为字节码
func (u *Unit) FromBool(b bool) error {
	if b {
		return u.FromStr("1")
	}
	return u.FromStr("0")
}

// ToBool 将字节码转为布尔值
func (u *Unit) ToBool() (bool, error) {
	return u.String() == "1", nil
}

// FromTime 将时间转为字节码
func (u *Unit) FromTime(t time.Time) error {
	return u.FromInt32(int32(t.Unix()))
}

// ToTime 将字节码转为时间
func (u *Unit) ToTime() (t time.Time, err error) {
	var n int32
	if n, err = u.ToInt32(); err == nil {
		t = time.Unix(int64(n), 0)
	}
	return
}
