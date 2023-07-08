package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

// GetColumnChanges 只保留匹配字段标签的数据
func GetColumnChanges(cols []string, changes map[string]any) map[string]any {
	result := make(map[string]any)
	for _, c := range cols {
		if val, ok := changes[c]; ok {
			result[c] = val
		}
	}
	return result
}

// GetStructTags Read all tags in a object
func GetStructTags(v any) map[string]reflect.StructTag {
	tags := make(map[string]reflect.StructTag)
	vt := GetIndirectType(v)
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		tags[field.Name] = field.Tag
	}
	return tags
}

// StructBuilder 已经解析的标签
type StructBuilder struct {
	FieldNames []string
	FieldTypes map[string]reflect.Type
	tags       map[string]*Tagger
}

// NewStructBuilder Read all tags in a object
func NewStructBuilder(v any) *StructBuilder {
	t := &StructBuilder{tags: make(map[string]*Tagger)}
	vt := GetIndirectType(v)
	if vt.Kind() != reflect.Struct {
		return t
	}
	num := vt.NumField()
	t.FieldNames = make([]string, num, num)
	for i := 0; i < num; i++ {
		field := vt.Field(i)
		t.FieldNames[i] = field.Name
		t.FieldTypes[field.Name] = field.Type
		t.tags[field.Name] = ParseTag(field.Tag)
	}
	return t
}

// GetFieldTag Read tag for a field and key
func (t *StructBuilder) GetFieldTag(name, key string) string {
	if tag, ok := t.tags[name]; ok {
		return tag.Get(key)
	}
	return ""
}

// GetColumnTags 读取结构体的字段标签
func (t *StructBuilder) GetColumnTags(cols []string, key string, alias ...string) []string {
	for name, tag := range t.tags {
		tagVal := tag.Get(key)
		if tagVal == "" || tagVal == "-" {
			continue
		} else if strings.HasSuffix(tagVal, "inline") {
			sub := NewStructBuilder(t.FieldTypes[name])
			cols = sub.GetColumnTags(cols, key, alias...)
		} else {
			if len(alias) > 0 {
				tagVal = fmt.Sprintf("%s.%s", alias[0], tagVal)
			}
			cols = append(cols, tagVal)
		}
	}
	return cols
}
