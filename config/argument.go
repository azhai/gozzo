package config

import (
	"flag"
)

type ArgList struct {
	args    map[string]int
	Convert func(string) string
}

func NewArgList(conv func(string) string) *ArgList {
	return &ArgList{
		args:    make(map[string]int),
		Convert: conv,
	}
}

// Add 增加一些参数
func (t *ArgList) Add(args []string, uniq bool) int {
	for _, arg := range args {
		if t.Convert != nil {
			arg = t.Convert(arg)
		}
		if uniq { // 去重不需要计数
			t.args[arg] = 1
		} else if val, ok := t.args[arg]; ok {
			t.args[arg] = val + 1
		} else {
			t.args[arg] = 1
		}
	}
	return t.Size()
}

// Count 获得此参数计数
func (t *ArgList) Count(arg string) int {
	if len(t.args) == 0 {
		return 0
	}
	if val, ok := t.args[arg]; ok {
		return val
	}
	return 0
}

// Has 是否含有此参数
func (t *ArgList) Has(arg string) bool {
	return t.Count(arg) > 0
}

// Size 参数元素个数
func (t *ArgList) Size() int {
	return len(t.args)
}

// ReadArgs 读取命令行参数，不包括命名参数，且必须将命名参数放在前面
func ReadArgs(uniq bool, conv func(string) string) *ArgList {
	lst := NewArgList(conv)
	if flag.NArg() > 0 {
		lst.Add(flag.Args(), uniq)
	}
	return lst
}
