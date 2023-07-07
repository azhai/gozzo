package mapper_test

import (
	"testing"

	"github.com/azhai/gozzo/mapper"
	"github.com/stretchr/testify/assert"
)

// 连接配置
type ConnParams struct {
	Host     string         `json:"host" yaml:"host" toml:"host"`
	Port     int            `json:"port,omitempty" yaml:"port,omitempty" toml:"port"`
	Username string         `yaml/json:"username,omitempty" toml:"username"`
	Password string         `toml/yaml/json:"password"`
	Database string         `toml/yaml/json:"database"`
	Options  map[string]any `yaml/json:"options,omitempty" toml:"options"`
}

// go test -run=Burnish
func Test_01_Itag_Burnish(t *testing.T) {
	tagDict := mapper.GetStructTags(ConnParams{})
	tag := mapper.BurnishTag(tagDict["Port"], "json", "toml", "yaml")
	assert.Equal(t, "port,omitempty", tag.Get("yaml"))
	assert.Equal(t, `json:"port,omitempty" toml:"port" yaml:"port,omitempty"`, tag.String())
}

// go test -run=Parse
func Test_02_Itag_Parse(t *testing.T) {
	tagDict := mapper.GetStructTags(ConnParams{})
	tag := mapper.ParseTag(tagDict["Username"])
	assert.Equal(t, "username,omitempty", tag.Get("yaml"))
	assert.Equal(t, `json:"username,omitempty" toml:"username" yaml:"username,omitempty"`, tag.String())
	tag = mapper.ParseTag(tagDict["Password"])
	assert.Equal(t, "password", tag.Get("yaml"))
	assert.Equal(t, `json:"password" toml:"password" yaml:"password"`, tag.String())
}
