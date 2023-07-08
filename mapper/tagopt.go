package mapper

import "strings"

// TagOpt 标签选项
type TagOpt struct {
	Name      string
	Convert   string
	OmitEmpty bool
	Inline    bool
}

// NewTagOpt 解析标签内容
func NewTagOpt(name, tag string) *TagOpt {
	tag = strings.TrimSpace(tag)
	if tag == "-" || tag == "" && name == "" {
		return nil
	}
	opt := &TagOpt{Name: name}
	if tag == "" {
		return opt
	} else if tag == "extends" {
		opt.Inline = true
		return opt
	}
	var found bool
	if opt.Name, tag, found = strings.Cut(tag, ","); found {
		_ = opt.ParseOption(tag)
	}
	return opt
}

// ParseOption 解析可选项
func (p *TagOpt) ParseOption(tag string) error {
	var head string
	tag = strings.ToLower(tag)
	for tag != "" {
		tag = strings.TrimLeft(tag, " \t\r\n")
		head, tag, _ = strings.Cut(tag, ",")
		if head == "omitempty" {
			p.OmitEmpty = true
		} else if head == "inline" || head == "extends" {
			p.Inline = true
		} else {
			p.Convert = head
		}
	}
	return nil
}
