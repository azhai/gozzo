package filesystem

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/azhai/gozzo/match"
)

// MapStrList 将函数应用到每个元素
func MapStrList(data []string, mf func(string) string, ff func(string) bool) []string {
	var result []string
	for _, v := range data {
		if ff != nil && ff(v) == false {
			continue
		}
		if mf != nil {
			v = mf(v)
		}
		result = append(result, v)
	}
	return result
}

// FindFiles 遍历目录下的文件，递归方法
func FindFiles(dir, ext string, excls ...string) (map[string]os.FileInfo, error) {
	result := make(map[string]os.FileInfo)
	exclMatchers := match.NewGlobs(MapStrList(excls, func(s string) string {
		if strings.HasSuffix(s, string(filepath.Separator)) {
			return s + "*" // 匹配所有目录下所有文件和子目录
		}
		return s
	}, nil))
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil { // 终止
			return err
		} else if exclMatchers.MatchAny(path, false) { // 跳过
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if info.Mode().IsRegular() {
			if ext == "" || strings.HasSuffix(info.Name(), ext) {
				result[path] = info
			}
		}
		return nil
	})
	return result, err
}

// CopyFiles 复制文件到新目录下
func CopyFiles(dest, src string, files map[string]string, force bool) (err error) {
	var content []byte
	for filename, toname := range files {
		if toname == "" {
			toname = filename
		}
		destFile := filepath.Join(dest, toname)
		if !force && File(destFile).IsExist() { // 不要覆盖
			continue
		}
		srcFile := filepath.Join(src, filename)
		if content, err = os.ReadFile(srcFile); err != nil {
			return
		}
		err = os.WriteFile(destFile, content, DefaultFileMode)
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
