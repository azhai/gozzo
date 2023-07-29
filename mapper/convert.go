package mapper

import (
	"net/url"
	"reflect"
)

// ConvertField 字段转换函数
type ConvertField func(field *StructField, opt *TagOpt) error

// TravelStruct 遍历和转换结构体的字段
func TravelStruct(obj any, tag string, omit bool, convert ConvertField) (err error) {
	var builder *StructBuilder
	if builder = NewStructBuilder(obj); builder == nil {
		return
	}
	fields, opts := builder.GetFieldTagOpts(tag, NameLowerFirst)
	for i, opt := range opts {
		if opt == nil { // 使用"-"标签跳过
			continue
		}
		if omit && opt.OmitEmpty {
			if fields[i].IsEmptyValue() { // omitEmpty时跳过为空的字段
				continue
			}
		}
		err = convert(fields[i], opt)
	}
	return
}

// DecodeToDict 将结构体转为哈希表
func DecodeToDict(obj any) (Dict, error) {
	data := make(Dict)
	err := TravelStruct(obj, "json", true,
		func(field *StructField, opt *TagOpt) (err error) {
			if field.Type.Kind() == reflect.Struct {
				data[opt.Name], err = DecodeToDict(field.Value)
			} else {
				data[opt.Name] = field.Value.Interface()
			}
			return
		})
	return data, err
}

// EncodeFromDict 将哈希表转为结构体
func EncodeFromDict(data Dict, obj any) error {
	return TravelStruct(obj, "json", false,
		func(field *StructField, opt *TagOpt) (err error) {
			if val, ok := data[opt.Name]; ok {
				err = field.SetValue(val, opt.ConvType)
			}
			return
		})
}

// EncodeFromURL 将url参数转为结构体
func EncodeFromURL(data url.Values, obj any) error {
	return TravelStruct(obj, "json", false,
		func(field *StructField, opt *TagOpt) (err error) {
			if val := data.Get(opt.Name); val != "" {
				err = field.SetString(val)
			}
			return
		})
}

// SetDictValue Change the value of dict field
// func SetDictValue(field *StructField, target Dict, name string) (Dict, error) {
// 	var err error
// 	switch field.Type.Kind() {
// 	case reflect.String:
// 		target[name] = field.Value.String()
// 	case reflect.Int:
// 		target[name] = int(field.Value.Int())
// 	case reflect.Int8:
// 		target[name] = int8(field.Value.Int())
// 	case reflect.Int16:
// 		target[name] = int16(field.Value.Int())
// 	case reflect.Int32:
// 		target[name] = int32(field.Value.Int())
// 	case reflect.Int64:
// 		target[name] = field.Value.Int()
// 	case reflect.Uint:
// 		target[name] = int(field.Value.Uint())
// 	case reflect.Uint8:
// 		target[name] = int8(field.Value.Uint())
// 	case reflect.Uint16:
// 		target[name] = int16(field.Value.Uint())
// 	case reflect.Uint32:
// 		target[name] = int32(field.Value.Uint())
// 	case reflect.Uint64:
// 		target[name] = field.Value.Uint()
// 	case reflect.Float32:
// 		target[name] = float32(field.Value.Float())
// 	case reflect.Float64:
// 		target[name] = field.Value.Float()
// 	case reflect.Bool:
// 		target[name] = field.Value.Bool()
// 	case reflect.Slice:
// 		vt := datatype.GetIndirectType(field.Value.Type())
// 		val := reflect.MakeSlice(vt.Elem(), 0, 0)
// 		target[name] = val.Interface()
// 	case reflect.Map:
// 		vt := datatype.GetIndirectType(field.Value.Type())
// 		mt := reflect.MapOf(vt.Key(), vt.Elem())
// 		val := reflect.MakeMap(mt)
// 		target[name] = val.Interface()
// 	case reflect.Struct:
// 		target[name], err = DecodeToDict(field.Value)
// 	}
// 	return target, err
// }
