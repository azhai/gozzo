package mapper

import (
	"reflect"
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
	Name    string
	Value   reflect.Value
	Type    reflect.Type
	Tag     *Tagger
	TagOpts map[string]*TagOpt
}

// GetTag Read tag for a key
func (f *StructField) GetTag(key string) string {
	return f.Tag.Get(key)
}

// GetTagOpt Read and parse tag for a key
func (f *StructField) GetTagOpt(key string) (opt *TagOpt) {
	var ok bool
	if opt, ok = f.TagOpts[key]; !ok {
		opt = NewTagOpt(f.Name, f.GetTag(key), NameNoChange)
		f.TagOpts[key] = opt
	}
	return
}

// SetValue Change the value of field
func (f *StructField) SetValue(val any) (err error) {
	vx := reflect.ValueOf(val)
	if vx.Kind() == f.Type.Kind() {
		f.Value.Set(vx)
		return
	}
	// TODO: 类型不同
	return
}

// SetDict Change the value of dict field
func (f *StructField) SetDict(name string, target Dict) (err error) {
	switch f.Type.Kind() {
	case reflect.String:
		target[name] = f.Value.String()
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		target[name] = f.Value.Int()
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		target[name] = f.Value.Uint()
	case reflect.Float64, reflect.Float32:
		target[name] = f.Value.Float()
	case reflect.Bool:
		target[name] = f.Value.Bool()
		// TODO: 补充其他类型，如 reflect.Slice
	}
	return
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
			Name:    field.Name,
			Value:   vv.FieldByIndex([]int{i}),
			Type:    field.Type,
			Tag:     ParseTag(field.Tag),
			TagOpts: make(map[string]*TagOpt),
		}
	}
	return builder
}

// GetFieldTagOpts 读取结构体的字段标签
func (b *StructBuilder) GetFieldTagOpts(key string,
	fields []*StructField, opts []*TagOpt,
) ([]*StructField, []*TagOpt) {
	for _, name := range b.Names {
		field := b.Fields[name]
		opt := field.GetTagOpt(key)
		if opt == nil {
			continue
		} else if opt.Inline {
			sub := NewStructBuilder(field.Value)
			fields, opts = sub.GetFieldTagOpts(key, fields, opts)
		} else {
			fields = append(fields, field)
			opts = append(opts, opt)
		}
	}
	return fields, opts
}
