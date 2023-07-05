//go:build windows

package logging

import (
	"os"
	"regexp"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}

// ignoreWinDisk 忽略Windows盘符
func ignoreWinDisk(absPath string) string {
	// 解决windows下zap不能识别路径中的盘符问题
	re := regexp.MustCompile(`^[A-Za-z]:`)
	return re.ReplaceAllLiteralString(absPath, "")
}
