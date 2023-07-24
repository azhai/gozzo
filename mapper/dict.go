package mapper

import (
	"fmt"
)

type Dict = map[string]any

func EncodeFrom(v any) (d Dict, err error) {
	d = make(map[string]any)
	var builder *StructBuilder
	if builder = NewStructBuilder(v); builder == nil {
		return
	}
	for _, opt := range builder.GetTagOpts("json") {
		fmt.Println(opt)
		field := builder.Fields[opt.FieldName]
		d[opt.Name] = field.Value.Int()
	}
	return
}

func DecodeTo(d Dict, v any) (err error) {
	var builder *StructBuilder
	if builder = NewStructBuilder(v); builder == nil {
		return
	}
	for _, opt := range builder.GetTagOpts("json") {
		if val, ok := d[opt.Name]; ok {
			field := builder.Fields[opt.FieldName]
			field.Value.SetInt(val.(int64))
		}
	}
	return
}
