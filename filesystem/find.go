package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/azhai/gozzo/match"
)

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
