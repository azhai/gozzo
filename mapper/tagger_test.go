package mapper_test

import (
	"reflect"
	"testing"

	"github.com/azhai/gozzo/mapper"
	"github.com/stretchr/testify/assert"
)

// getStructTags Read all tags in a object
func getStructTags(v any) map[string]reflect.StructTag {
	tags := make(map[string]reflect.StructTag)
	vt := mapper.GetIndirectType(v)
	if vt.Kind() != reflect.Struct {
		return tags
	}
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		tags[field.Name] = field.Tag
	}
	return tags
}

// ConnParams 连接配置
type ConnParams struct {
	Host     string         `json:"host" yaml:"host" toml:"host"`
	Port     int            `json:"port,omitempty" yaml:"port,omitempty" toml:"port"`
	Username string         `yaml/json:"username" toml:"username"`
	Password string         `toml/yaml/json:"password"`
	Database string         `toml/yaml/json:"database"`
	Options  map[string]any `yaml/json:"options,omitempty" toml:"options"`
}

// go test -run=TagBurnish
func Test01_TagBurnish(t *testing.T) {
	tagDict := getStructTags(ConnParams{})
	tag := mapper.BurnishTag(tagDict["Port"], "json", "toml", "yaml")
	assert.Equal(t, "port,omitempty", tag.Get("yaml"))
	assert.Equal(t, `json:"port,omitempty" toml:"port" yaml:"port,omitempty"`, tag.String())
}

// go test -run=TagParse
func Test02_TagParse(t *testing.T) {
	tagDict := getStructTags(ConnParams{})
	tag := mapper.ParseTag(tagDict["Password"])
	assert.Equal(t, "password", tag.Get("yaml"))
	assert.Equal(t, `json:"password" toml:"password" yaml:"password"`, tag.String())
	tag = mapper.ParseTag(tagDict["Options"])
	assert.Equal(t, "options,omitempty", tag.Get("yaml"))
	assert.Equal(t, `json:"options,omitempty" toml:"options" yaml:"options,omitempty"`, tag.String())
}
