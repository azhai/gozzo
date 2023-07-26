package mapper_test

import (
	"testing"

	"github.com/azhai/gozzo/mapper"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

// go test -run=DictDecode
func Test11_DictDecode(t *testing.T) {
	params := &ConnParams{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Database: "test",
		// Options:  make(map[string]any),
	}
	dict, err := mapper.DecodeToDict(params)
	pp.Println(dict)
	assert.NoError(t, err)
	assert.Equal(t, params.Host, dict["host"])
	assert.Equal(t, params.Port, dict["port"])
	assert.Equal(t, params.Username, dict["username"])
	assert.Equal(t, params.Password, dict["password"])
	assert.Equal(t, params.Database, dict["database"])
	assert.Nil(t, dict["options"])
}

// go test -run=DictEncode
func Test12_DictEncode(t *testing.T) {
	dict := map[string]any{
		"host":     "127.0.0.1",
		"port":     3306,
		"username": "root",
		"database": "test",
		"options":  make(map[string]any),
	}
	params := new(ConnParams)
	err := mapper.EncodeFromDict(dict, params)
	pp.Println(params)
	assert.NoError(t, err)
	assert.Equal(t, params.Host, dict["host"])
	assert.Equal(t, params.Port, dict["port"])
	assert.Equal(t, params.Username, dict["username"])
	assert.Equal(t, params.Password, "")
	assert.Equal(t, params.Database, dict["database"])
	assert.Equal(t, params.Options, dict["options"])
}
