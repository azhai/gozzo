package mapper

// DecodeToDict 将结构体转为哈希表
func DecodeToDict(v any) (d Dict, err error) {
	var builder *StructBuilder
	if builder = NewStructBuilder(v); builder == nil {
		return
	}
	d = make(Dict)
	fields, opts := builder.GetFieldTagOpts("json", NameLowerFirst)
	for i, opt := range opts {
		if opt == nil { // 使用"-"标签跳过
			continue
		}
		d, err = fields[i].SetDict(d, opt.Name, opt.OmitEmpty)
	}
	return
}

type Decoder struct{}

func (r *Decoder) Decode() error {
	return nil
}

func (r *Decoder) DecodeSingle() error {
	return nil
}

func (r *Decoder) DecodeSlice() error {
	return nil
}

func (r *Decoder) DecodeStruct() error {
	return nil
}

func (r *Decoder) DecodeCustom() error {
	return nil
}
