package mapper

import (
	"reflect"
	"strconv"
)

// Dict 哈希表
type Dict = map[string]any

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
func (f *StructField) GetTagOpt(key string, caser NameCase) (opt *TagOpt) {
	var ok bool
	if opt, ok = f.TagOpts[key]; !ok {
		opt = NewTagOpt(f.Name, f.GetTag(key), caser)
		f.TagOpts[key] = opt
	}
	return
}

// SetValue Change the value of field
func (f *StructField) SetValue(val any, conv string) (err error) {
	vx := reflect.ValueOf(val)
	if vx.Kind() == f.Type.Kind() {
		f.Value.Set(vx)
		return
	} else if conv == "" {
		return
	}
	// 类型不同且有转换函数
	switch vx.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := vx.Int()
		if conv == "string" {
			f.Value.SetString(strconv.FormatInt(value, 10))
		} else if conv == "int" {
			f.Value.SetInt(value)
		} else if conv == "uint" {
			f.Value.SetUint(uint64(value))
		} else if conv == "float" {
			f.Value.SetFloat(float64(value))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := vx.Uint()
		if conv == "string" {
			f.Value.SetString(strconv.FormatUint(value, 10))
		} else if conv == "int" {
			f.Value.SetInt(int64(value))
		} else if conv == "uint" {
			f.Value.SetUint(value)
		} else if conv == "float" {
			f.Value.SetFloat(float64(value))
		}
	case reflect.Float32, reflect.Float64:
		value := vx.Float()
		if conv == "string" {
			f.Value.SetString(strconv.FormatFloat(value, 'G', -1, 64))
		} else if conv == "int" {
			f.Value.SetInt(int64(value))
		} else if conv == "uint" {
			f.Value.SetUint(uint64(value))
		} else if conv == "float" {
			f.Value.SetFloat(value)
		}
	}
	return
}

// SetDict Change the value of dict field
func (f *StructField) SetDict(target Dict, name string, omitEmpty bool) (Dict, error) {
	var err error
	kind := f.Type.Kind()
	if omitEmpty {
		if kind == reflect.Pointer && f.Value.IsNil() {
			return target, err
		}
		if kind != reflect.Pointer && f.Value.IsZero() {
			return target, err
		}
	}
	switch kind {
	case reflect.String:
		target[name] = f.Value.String()
	case reflect.Int:
		target[name] = int(f.Value.Int())
	case reflect.Int8:
		target[name] = int8(f.Value.Int())
	case reflect.Int16:
		target[name] = int16(f.Value.Int())
	case reflect.Int32:
		target[name] = int32(f.Value.Int())
	case reflect.Int64:
		target[name] = f.Value.Int()
	case reflect.Uint:
		target[name] = int(f.Value.Uint())
	case reflect.Uint8:
		target[name] = int8(f.Value.Uint())
	case reflect.Uint16:
		target[name] = int16(f.Value.Uint())
	case reflect.Uint32:
		target[name] = int32(f.Value.Uint())
	case reflect.Uint64:
		target[name] = f.Value.Uint()
	case reflect.Float32:
		target[name] = float32(f.Value.Float())
	case reflect.Float64:
		target[name] = f.Value.Float()
	case reflect.Bool:
		target[name] = f.Value.Bool()
	case reflect.Slice:
		vt := GetIndirectType(f.Value.Type())
		val := reflect.MakeSlice(vt.Elem(), 0, 0)
		target[name] = val.Interface()
	case reflect.Map:
		vt := GetIndirectType(f.Value.Type())
		mt := reflect.MapOf(vt.Key(), vt.Elem())
		val := reflect.MakeMap(mt)
		target[name] = val.Interface()
	case reflect.Struct:
		target[name], err = DecodeToDict(f.Value)
	}
	return target, err
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
	vv := reflect.Indirect(reflect.ValueOf(v))
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

// getFieldTagOpts 读取结构体的字段标签
func (b *StructBuilder) getFieldTagOpts(key string, caser NameCase,
	fields []*StructField, opts []*TagOpt,
) ([]*StructField, []*TagOpt) {
	for _, name := range b.Names {
		field := b.Fields[name]
		opt := field.GetTagOpt(key, caser)
		if opt == nil {
			continue
		} else if opt.Inline {
			sub := NewStructBuilder(field.Value)
			fields, opts = sub.getFieldTagOpts(key, caser, fields, opts)
		} else {
			fields = append(fields, field)
			opts = append(opts, opt)
		}
	}
	return fields, opts
}

// GetFieldTagOpts 读取结构体的字段标签
func (b *StructBuilder) GetFieldTagOpts(key string, caser NameCase,
) ([]*StructField, []*TagOpt) {
	return b.getFieldTagOpts(key, caser, nil, nil)
}
