package cryptogy

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/muyo/sno"
)

var (
	ciper      = Cipher()
	saltPasswd ICipher
)

type ICipher interface {
	CreatePassword(plainText string) string
	VerifyPassword(plainText, cipherText string) bool
}

func Cipher() ICipher {
	if saltPasswd == nil { // 8位salt值，用$符号分隔开
		saltPasswd = NewSaltPassword(8, "$")
	}
	return saltPasswd
}

func NewSerialNo(n byte) string {
	return sno.New(n).String()
}

func NewTimeSerialNo(n byte, t time.Time) string {
	return sno.NewWithTime(n, t).String()
}

func Md5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	cipher := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipher)
}

func VerifyPassword(plainText, cipherText string) bool {
	return ciper.VerifyPassword(plainText, cipherText)
}

func CreatePassword(password string) string {
	return ciper.CreatePassword(password)
}

func CreateMd5Password(password string) string {
	password = Md5(strings.Repeat(password, 2)) // 客户端进行了哈希
	return CreatePassword(password)
}

func GetPasswordChanges(password string) map[string]any {
	return map[string]any{
		"password": CreateMd5Password(password),
	}
}

func CreateToken(prefix []byte, tailsize int) string {
	tailno := RandSalt(tailsize * 2)
	return hex.EncodeToString(prefix) + tailno
}
