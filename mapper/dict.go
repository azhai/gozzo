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
	fields, opts := builder.GetFieldTagOpts("json", nil, nil)
	for i, opt := range opts {
		fmt.Println(opt)
		d[opt.Name] = fields[i].Value.Int()
	}
	return
}

func DecodeTo(d Dict, v any) (err error) {
	var builder *StructBuilder
	if builder = NewStructBuilder(v); builder == nil {
		return
	}
	fields, opts := builder.GetFieldTagOpts("json", nil, nil)
	for i, opt := range opts {
		if val, ok := d[opt.Name]; ok {
			fields[i].Value.SetInt(val.(int64))
		}
	}
	return
}
