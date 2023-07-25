package mapper

import (
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	NameNoChange NameCase = iota
	NameLowerFirst
	NameLowerCase
	NameCamelCase
	NameSnakeCase
)

type NameCase int

func ToCase(name string, caser NameCase) string {
	switch caser {
	default:
		return name
	case NameLowerFirst:
		return strcase.ToLowerCamel(name)
	case NameLowerCase:
		return strings.ToLower(name)
	case NameCamelCase:
		return strcase.ToCamel(name)
	case NameSnakeCase:
		return strcase.ToSnake(name)
	}
}

// TagOpt 标签选项
type TagOpt struct {
	FieldName string
	Name      string
	Convert   string
	OmitEmpty bool
	Inline    bool
}

// NewTagOpt 解析标签内容
func NewTagOpt(name, tag string, caser NameCase) *TagOpt {
	tag = strings.TrimSpace(tag)
	if tag == "-" || tag == "" && name == "" {
		return nil
	}
	opt := &TagOpt{FieldName: name, Name: ToCase(name, caser)}
	if tag == "" {
		return opt
	}
	head, tail, found := strings.Cut(tag, ",")
	if head == "extends" {
		opt.Inline = true
	} else if head != "" {
		opt.Name = head
	}
	if found {
		_ = opt.ParseOption(tail)
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
