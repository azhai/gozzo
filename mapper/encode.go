package mapper

// EncodeFromDict 将哈希表转为结构体
func EncodeFromDict(d Dict, v any) (err error) {
	var builder *StructBuilder
	if builder = NewStructBuilder(v); builder == nil {
		return
	}
	fields, opts := builder.GetFieldTagOpts("json", NameLowerFirst)
	for i, opt := range opts {
		if opt == nil { // 使用"-"标签跳过
			continue
		}
		if val, ok := d[opt.Name]; ok {
			err = fields[i].SetValue(val, opt.Convert)
		}
	}
	return
}
