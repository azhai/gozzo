package logging

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var callerSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	callerSourceDir = sourceDir(file)
}

func sourceDir(file string) string {
	dir := filepath.Dir(filepath.Dir(file))
	s := filepath.Dir(dir)
	if filepath.Base(s) != "gorm.io" {
		s = dir
	}
	return filepath.ToSlash(s) + "/"
}

// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, callerSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
