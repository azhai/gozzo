package cryptogy

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// RandSalt 产生随机salt
func RandSalt(size int) string {
	buf := make([]byte, (size+1)/2)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)[:size]
	}
	return ""
}

// SaltPassword 带salt值的sha256密码哈希
type SaltPassword struct {
	saltLen int
	saltSep string
	*MacHash
}

func NewSaltPassword(len int, sep string) *SaltPassword {
	return &SaltPassword{
		saltLen: len, saltSep: sep,
		MacHash: NewMacHash(sha256.New),
	}
}

// CreatePassword 设置密码
func (p *SaltPassword) CreatePassword(plainText string) string {
	saltValue := RandSalt(p.saltLen)
	cipherText := p.SetKey(saltValue).Sign(plainText)
	return saltValue + p.saltSep + cipherText
}

// VerifyPassword 校验密码
func (p *SaltPassword) VerifyPassword(plainText, cipherText string) bool {
	pieces := strings.SplitN(cipherText, p.saltSep, 2)
	if len(pieces) == 2 {
		return p.SetKey(pieces[0]).Verify(plainText, pieces[1])
	}
	return false
}
