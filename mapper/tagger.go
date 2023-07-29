package mapper

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/azhai/gozzo/datatype"
	"github.com/azhai/gozzo/match"
)

// TagKeySep tag分隔符，约定大于配置
const TagKeySep = "/"

// Tagger Yet another StructTag
type Tagger struct {
	alias   map[string]string // 例如{"db":"json", "yaml":"json"}
	data    map[string]string // 例如{"json":"port,omitempty"}
	changed bool
	lock    sync.RWMutex
	reflect.StructTag
}

// NewTagger 创建空的tag
func NewTagger() *Tagger {
	t := new(Tagger)
	t.alias = make(map[string]string)
	t.data = make(map[string]string)
	return t
}

// ParseTag 解析全部tags
func ParseTag(tag reflect.StructTag) *Tagger {
	t := NewTagger()
	t.StructTag, t.changed = tag, true
	word, key, chunk := match.Word(tag), "", ""
	for word != "" {
		if word = word.SkipAnyChar(" "); word == "" {
			break
		}
		if key, word = word.MatchSubString(":\""); key == "" {
			break
		}
		chunk, word = word.MatchSubString("\"")
		value, err := strconv.Unquote("\"" + chunk + "\"")
		if err != nil {
			break
		}
		t.Append(key, value)
	}
	return t
}

// BurnishTag 保留部分tags
func BurnishTag(tag reflect.StructTag, names ...string) *Tagger {
	t := NewTagger()
	t.StructTag, t.changed = tag, true
	for _, key := range names {
		if value, ok := tag.Lookup(key); ok {
			t.Append(key, value)
		}
	}
	return t
}

// Build 组装为字节数组
func (t *Tagger) Build(data []byte, name, value string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, byte(':'), byte('"'))
	data = append(data, []byte(value)...)
	data = append(data, byte('"'), byte(' '))
	return data
}

// String 转为字符串格式
func (t *Tagger) String() string {
	if !t.changed {
		return string(t.StructTag)
	}
	var data []byte
	if len(t.alias) > len(t.data) {
		for _, key := range datatype.SortedMapKeys(t.alias) {
			if value := t.Get(key); value != "" {
				data = t.Build(data, key, value)
			}
		}
	} else {
		for _, name := range datatype.SortedMapKeys(t.data) {
			data = t.Build(data, name, t.data[name])
		}
	}
	result := strings.TrimSpace(string(data))
	t.StructTag = reflect.StructTag(result)
	t.changed = false
	return result
}

// Get returns the value associated with key in the tag string
func (t *Tagger) Get(key string) string {
	v, _ := t.Lookup(key)
	return v
}

// Lookup Returns a tag from the tag data
func (t *Tagger) Lookup(key string) (string, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if name, ok := t.alias[key]; ok {
		key = name
	}
	value, ok := t.data[key]
	return value, ok
}

// Append Sets a tag in the tag data map
func (t *Tagger) Append(key, value string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	names := strings.Split(key, TagKeySep)
	name := names[len(names)-1]
	for _, key = range names {
		t.alias[key] = name
	}
	t.changed, t.data[name] = true, value
}

// Delete Deletes a tag
func (t *Tagger) Delete(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	delete(t.alias, key)
	delete(t.data, key)
	t.changed = true
}
