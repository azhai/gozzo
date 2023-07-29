package datatype

import (
	"strconv"
	"time"
)

// Unit 数据转换单元
type Unit []byte

// String 将字节码转为字符串
func (u Unit) String() string {
	return string(u)
}

// FromStr 将字符串转为字节码
func (u *Unit) FromStr(s string) error {
	*u = []byte(s)
	return nil
}

// ToStr 将字节码转为字符串
func (u Unit) ToStr() (string, error) {
	return u.String(), nil
}

// FromInt 将整数转为字节码
func (u Unit) FromInt(n int) error {
	return u.FromStr(strconv.Itoa(n))
}

// ToInt 将字节码转为整数
func (u Unit) ToInt() (int, error) {
	return strconv.Atoi(u.String())
}

// FromInt64 将整数转为字节码
func (u Unit) FromInt64(n int64) error {
	return u.FromStr(strconv.FormatInt(n, 10))
}

// ToInt64 将字节码转为整数
func (u Unit) ToInt64() (int64, error) {
	return strconv.ParseInt(u.String(), 10, 64)
}

// MustInt64 将字节码转为整数
func (u Unit) MustInt64() int64 {
	n, err := u.ToInt64()
	if err != nil {
		panic(err)
	}
	return n
}

// FromInt32 将整数转为字节码
func (u Unit) FromInt32(n int32) error {
	return u.FromInt64(int64(n))
}

// ToInt32 将字节码转为整数
func (u Unit) ToInt32() (int32, error) {
	n, err := u.ToInt64()
	return int32(n), err
}

// FromInt16 将整数转为字节码
func (u Unit) FromInt16(n int16) error {
	return u.FromInt64(int64(n))
}

// ToInt16 将字节码转为整数
func (u Unit) ToInt16() (int16, error) {
	n, err := u.ToInt64()
	return int16(n), err
}

// FromInt8 将整数转为字节码
func (u Unit) FromInt8(n int8) error {
	return u.FromInt64(int64(n))
}

// ToInt8 将字节码转为整数
func (u Unit) ToInt8() (int8, error) {
	n, err := u.ToInt64()
	return int8(n), err
}

// FromUint 将整数转为字节码
func (u Unit) FromUint(n uint) error {
	return u.FromStr(strconv.FormatUint(uint64(n), 10))
}

// ToUint 将字节码转为整数
func (u Unit) ToUint() (uint, error) {
	n, err := u.ToUint64()
	return uint(n), err
}

// FromUint64 将整数转为字节码
func (u Unit) FromUint64(n uint64) error {
	return u.FromStr(strconv.FormatUint(n, 10))
}

// ToUint64 将字节码转为整数
func (u Unit) ToUint64() (uint64, error) {
	return strconv.ParseUint(u.String(), 10, 64)
}

// FromUint32 将整数转为字节码
func (u Unit) FromUint32(n uint32) error {
	return u.FromUint64(uint64(n))
}

// ToUint32 将字节码转为整数
func (u Unit) ToUint32() (uint32, error) {
	n, err := u.ToUint64()
	return uint32(n), err
}

// FromUint16 将整数转为字节码
func (u Unit) FromUint16(n uint16) error {
	return u.FromUint64(uint64(n))
}

// ToUint16 将字节码转为整数
func (u Unit) ToUint16() (uint16, error) {
	n, err := u.ToUint64()
	return uint16(n), err
}

// FromUint8 将整数转为字节码
func (u Unit) FromUint8(n uint8) error {
	return u.FromUint64(uint64(n))
}

// ToUint8 将字节码转为整数
func (u Unit) ToUint8() (uint8, error) {
	n, err := u.ToUint64()
	return uint8(n), err
}

// FromFloat64 将浮点数转为字节码
func (u Unit) FromFloat64(f float64) error {
	return u.FromStr(strconv.FormatFloat(f, 'G', -1, 64))
}

// ToFloat64 将字节码转为浮点数
func (u Unit) ToFloat64() (float64, error) {
	return strconv.ParseFloat(u.String(), 64)
}

// FromFloat32 将浮点数转为字节码
func (u Unit) FromFloat32(f float32) error {
	return u.FromFloat64(float64(f))
}

// ToFloat32 将字节码转为浮点数
func (u Unit) ToFloat32() (float32, error) {
	f, err := u.ToFloat64()
	return float32(f), err
}

// FromBool 将布尔值转为字节码
func (u *Unit) FromBool(b bool) error {
	return u.FromStr(strconv.FormatBool(b))
}

// ToBool 将字节码转为布尔值
func (u Unit) ToBool() (bool, error) {
	return strconv.ParseBool(u.String())
}

// FromTime 将日期时间转为字节码
func (u Unit) FromTime(t time.Time) error {
	return u.FromStr(t.Format(time.DateTime))
}

// ToTime 将字节码转为日期时间
func (u Unit) ToTime() (t time.Time, err error) {
	return time.Parse(time.DateTime, u.String())
}

// FromStamp 将时间戳转为字节码
func (u Unit) FromStamp(t time.Time) error {
	return u.FromInt32(int32(t.Unix()))
}

// ToStamp 将字节码转为时间戳
func (u Unit) ToStamp() (t time.Time, err error) {
	var n int64
	if n, err = u.ToInt64(); err == nil {
		t = time.Unix(n, 0)
	}
	return
}
