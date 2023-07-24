package mapper

import (
	"reflect"
	"strings"
)

// GetColumnChanges 只保留匹配字段标签的数据
func GetColumnChanges(cols []string, changes map[string]any,
) map[string]any {
	result := make(map[string]any)
	for _, c := range cols {
		if val, ok := changes[c]; ok {
			result[c] = val
		}
	}
	return result
}

// StructField 字段信息
type StructField struct {
	Name  string
	Type  reflect.Type
	Value reflect.Value
	*Tagger
}

// StructBuilder 已经解析的标签
type StructBuilder struct {
	Names  []string
	Fields map[string]*StructField
}

// NewStructBuilder Read all tags in a object
func NewStructBuilder(v any) *StructBuilder {
	vt := GetIndirectType(v)
	if vt.Kind() != reflect.Struct {
		return nil
	}
	num := vt.NumField()
	builder := &StructBuilder{
		Names:  make([]string, num, num),
		Fields: make(map[string]*StructField, num),
	}
	vv := reflect.ValueOf(v)
	for i := 0; i < num; i++ {
		field := vt.Field(i)
		builder.Names[i] = field.Name
		builder.Fields[field.Name] = &StructField{
			Name:   field.Name,
			Type:   field.Type,
			Value:  vv.FieldByIndex([]int{i}),
			Tagger: ParseTag(field.Tag),
		}
	}
	return builder
}

// GetFieldTag Read tag for a field and key
func (b *StructBuilder) GetFieldTag(name, key string) string {
	if field, ok := b.Fields[name]; ok {
		return field.Get(key)
	}
	return ""
}

// getTagOpts 读取结构体的字段标签
func (b *StructBuilder) getTagOpts(key string, cols []string, opts []*TagOpt,
) ([]string, []*TagOpt) {
	for _, name := range b.Names {
		tagVal := b.Fields[name].Get(key)
		if tagVal == "-" {
			continue
		} else if strings.HasSuffix(tagVal, "inline") ||
			strings.HasSuffix(tagVal, "extends") {
			sub := NewStructBuilder(b.Fields[name].Type)
			cols, opts = sub.getTagOpts(key, cols, opts)
		} else {
			opt := NewTagOpt(name, tagVal)
			opts = append(opts, opt)
			cols = append(cols, opt.Name)
		}
	}
	return cols, opts
}

// GetTagOpts 读取结构体的字段标签
func (b *StructBuilder) GetTagOpts(key string) []*TagOpt {
	_, opts := b.getTagOpts(key, nil, nil)
	return opts
}
