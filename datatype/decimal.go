package datatype

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func RoundN(x float64, n int) float64 {
	multiple := math.Pow10(n)
	return math.Round(x*multiple) / multiple
}

// Decimal 高精度小数
// 例如 123.45 表示为 Decimal{Value:int64(12345), Precision:2}
type Decimal struct {
	Value     int64 // 扩大后成为整数
	Precision int   // 小数点后位数，限制15以内
}

// NewDecimal 使用方法：NewDecimal(123.45, 2)
func NewDecimal(value float64, prec int) *Decimal {
	d := &Decimal{}
	d.SetPrecision(prec)
	d.SetFloat(value, d.Precision)
	return d
}

// ParseDecimal 使用方法：ParseDecimal("123.45", 2)
func ParseDecimal(text string, prec int) *Decimal {
	d := &Decimal{}
	d.SetPrecision(prec)
	if idx := strings.Index(text, "."); idx >= 0 {
		size := d.Precision + idx + 1
		if paddings := size - len(text); paddings > 0 {
			zeros := strings.Repeat("0", paddings)
			text = text[:idx] + text[idx+1:] + zeros
		} else {
			text = text[:idx] + text[idx+1:size]
		}
	}
	d.Value, _ = strconv.ParseInt(text, 10, 64)
	return d
}

func (d *Decimal) HasFraction() bool {
	if d.Precision <= 0 {
		return false
	}
	base := int64(math.Pow10(d.Precision))
	return d.Value%base != 0
}

func (d *Decimal) GetFloat() float64 {
	return float64(d.Value) / math.Pow10(d.Precision)
}

func (d *Decimal) SetFloat(value float64, expand int) {
	if expand > 0 {
		value *= math.Pow10(expand)
	}
	d.Value = int64(math.Round(value))
}

func (d *Decimal) SetPrecision(prec int) {
	if prec >= 15 {
		d.Precision = 15
	} else if prec <= 0 {
		d.Precision = 0
	} else {
		d.Precision = prec
	}
}

func (d *Decimal) ChangePrecision(offset int) {
	oldPrec := d.Precision
	d.SetPrecision(d.Precision + offset)
	offset = d.Precision - oldPrec
	if offset > 0 {
		d.Value *= int64(math.Pow10(offset))
	} else if offset < 0 {
		d.SetFloat(float64(d.Value), 0-offset)
	}
}

// Format 保留小数点之后末尾的0
func (d *Decimal) Format() string {
	// 多加一个前置0，兼容无整数部分的情况
	size := int(d.Precision) + 1
	tpl := "%0" + strconv.Itoa(size) + "d"
	result := fmt.Sprintf(tpl, d.Value)
	if sep := len(result) + 1 - size; sep > 0 {
		result = result[:sep] + "." + result[sep:]
	}
	return result
}

// String 不保留小数点之后末尾的0
func (d *Decimal) String() string {
	result := d.Format()
	// 分开去除，否则会去掉整数部分末尾的0
	result = strings.TrimRight(result, "0")
	result = strings.TrimRight(result, ".")
	return result
}
