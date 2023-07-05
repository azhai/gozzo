package filesystem

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultDirMode  = 0o755
	DefaultFileMode = 0o644
)

// FileSize 检查文件是否存在及大小
// -1, false 不合法的路径
// 0, false 路径不存在
// -1, true 存在文件夹
// >=0, true 文件并存在
func FileSize(path string) (int64, bool) {
	if path == "" {
		return -1, false
	}
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return 0, false
	}
	var size = int64(-1)
	if info.IsDir() == false {
		size = info.Size()
	}
	return size, true
}

func CreateFile(path string) (fp *os.File, err error) {
	// create dirs if file not exists
	if dir := filepath.Dir(path); dir != "." {
		err = os.MkdirAll(dir, DefaultDirMode)
	}
	if err == nil {
		flag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
		fp, err = os.OpenFile(path, flag, DefaultFileMode)
	}
	return
}

func OpenFile(path string, readonly, append bool) (fp *os.File, size int64, err error) {
	var exists bool
	size, exists = FileSize(path)
	if size < 0 {
		err = fmt.Errorf("path is directory or illegal")
		return
	}
	if exists {
		flag := os.O_RDWR
		if readonly {
			flag = os.O_RDONLY
		} else if append {
			flag |= os.O_APPEND
		}
		fp, err = os.OpenFile(path, flag, DefaultFileMode)
	} else if readonly == false {
		fp, err = CreateFile(path)
	}
	return
}

// LineCount 使用 wc -l 计算有多少行
func LineCount(filename string) int {
	var err error
	filename, err = filepath.Abs(filename)
	if err != nil {
		return -1
	}
	var out []byte
	out, err = exec.Command("wc", "-l", filename).Output()
	if err != nil {
		return -1
	}
	num := 0
	col := strings.SplitN(string(out), " ", 2)[0]
	if num, err = strconv.Atoi(col); err != nil {
		return -1
	}
	return num
}

// MkdirForFile 为文件路径创建目录
func MkdirForFile(path string) int64 {
	size, exists := FileSize(path)
	if size < 0 {
		return size
	}
	if !exists {
		dir := filepath.Dir(path)
		_ = os.MkdirAll(dir, DefaultDirMode)
	}
	return size
}
