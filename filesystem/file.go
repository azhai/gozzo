package filesystem

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
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

type FileHandler struct {
	path    string
	err     error
	handler *os.File
	os.FileInfo
}

func NewFileHandler(path string) *FileHandler {
	fh := &FileHandler{path: path}
	_ = fh.Stat()
	return fh
}

func (f *FileHandler) Stat() error {
	f.FileInfo, f.err = os.Stat(f.path)
	return f.err
}

func (f *FileHandler) Error() error {
	return f.err
}

func (f *FileHandler) IsExist() bool {
	return f.err == nil || !os.IsNotExist(f.err)
}

func (f *FileHandler) IsAllow() bool {
	return f.err == nil || !os.IsPermission(f.err)
}

func (f *FileHandler) Close() error {
	return f.handler.Close()
}

func (f *FileHandler) Create() *os.File {
	if f.IsExist() || f.IsDir() {
		return nil
	}
	f.handler, f.err = CreateFile(f.path)
	return f.handler
}

func (f *FileHandler) Open(flag int) *os.File {
	if f.IsExist() && f.IsAllow() && !f.IsDir() {
		f.handler, f.err = OpenFile(f.path, flag)
	} else if flag&os.O_CREATE == os.O_CREATE {
		f.Create()
	}
	return f.handler
}

func (f *FileHandler) GetDims() (int, int) {
	if f.Open(os.O_RDONLY); f.err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", f.path, f.err)
		return 0, 0
	}
	img, _, err := image.DecodeConfig(f.handler)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", f.path, err)
		return 0, 0
	}
	return img.Width, img.Height
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

func OpenFile(path string, flag int) (fp *os.File, err error) {
	fp, err = os.OpenFile(path, flag, DefaultFileMode)
	if os.IsNotExist(err) {
		fp, err = CreateFile(path)
	}
	return
}

func WriteFile(path string, data []byte, append bool) error {
	flag := os.O_RDWR | os.O_TRUNC
	if append {
		flag = os.O_RDWR | os.O_APPEND
	}
	fp, err := OpenFile(path, flag)
	if err == nil {
		defer fp.Close()
		_, err = fp.Write(data)
	}
	return err
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
