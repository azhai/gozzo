package rewrite

import (
	"bytes"
	"go/format"
	"os"
	"strings"

	"github.com/azhai/gozzo/filesystem"
	"golang.org/x/tools/imports"
)

// TrimComment 去掉注释两边的空白
func TrimComment(c string) string {
	c = strings.TrimSpace(c)
	if strings.HasPrefix(c, "//") {
		c = strings.TrimSpace(c[2:])
	}
	return c
}

// FormatGolangCode 格式化代码，如果出错返回原内容
func FormatGolangCode(src []byte) ([]byte, error) {
	_src, err := format.Source(src)
	if err == nil {
		src = _src
	}
	return src, err
}

// SaveCodeToFile 将go代码保存到文件
func SaveCodeToFile(filename string, codeText []byte) ([]byte, error) {
	filesystem.MkdirForFile(filename)
	err := os.WriteFile(filename, codeText, filesystem.DefaultFileMode)
	return codeText, err
}

// PrettifyGolangFile 读出来go代码，重新写入文件
func PrettifyGolangFile(filename string, cleanImports bool) (changed bool, err error) {
	var srcCode, dstCode []byte
	if srcCode, err = os.ReadFile(filename); err != nil {
		return
	}
	dstCode, err = writeGolangFile(filename, srcCode, cleanImports)
	if bytes.Compare(srcCode, dstCode) != 0 {
		changed = true
	}
	return
}

// writeGolangFile 将代码整理后写入文件
func writeGolangFile(filename string, codeText []byte,
	cleanImports bool) (srcCode []byte, err error) {
	// Formart/Prettify the code 格式化代码
	if srcCode, err = FormatGolangCode(codeText); err != nil {
		srcCode = codeText
	}
	defer func() {
		_, errSave := SaveCodeToFile(filename, srcCode)
		if err == nil {
			err = errSave
		}
	}()
	if err == nil && cleanImports { // 清理 import
		cs := NewCodeSource()
		_ = cs.SetSource(srcCode)
		if cs.CleanImports() > 0 {
			srcCode, err = cs.GetContent()
		}
	}
	return
}

// splitImports 分组排序引用包
func splitImports(filename string, srcCode []byte, err error) ([]byte, error) {
	if err == nil {
		var dstCode []byte
		// Split the imports in two groups: go standard and the third parts
		dstCode, err = imports.Process(filename, srcCode, nil)
		if err == nil {
			return SaveCodeToFile(filename, dstCode)
		}
	}
	return srcCode, err
}

// WriteGolangFilePrettify 美化并输出go代码到文件
func WriteGolangFilePrettify(filename string, codeText []byte) ([]byte, error) {
	srcCode, err := writeGolangFile(filename, codeText, false)
	return splitImports(filename, srcCode, err)
}

// WriteGolangFileCleanImports 美化和整理导入，并输出go代码到文件
func WriteGolangFileCleanImports(filename string, codeText []byte) ([]byte, error) {
	srcCode, err := writeGolangFile(filename, codeText, true)
	return splitImports(filename, srcCode, err)
}

// WritePackage 将包中的Go文件格式化，如果提供了pkgname则用作新包名
func WritePackage(pkgPath, pkgName string) error {
	if pkgName != "" {
		// TODO: 替换包名
	}
	files, err := filesystem.FindFiles(pkgPath, ".go")
	if err != nil {
		return err
	}
	var content []byte
	for filename := range files {
		content, err = os.ReadFile(filename)
		if err != nil {
			break
		}
		_, err = WriteGolangFilePrettify(filename, content)
		if err != nil {
			break
		}
	}
	return err
}

// WriteWithImports 注入导入声明
func WriteWithImports(pkg string, source []byte,
	imports map[string]string) (*CodeSource, error) {
	cs := NewCodeSource()
	if err := cs.SetPackage(pkg); err != nil {
		return cs, err
	}
	// 添加可能引用的包，后面再尝试删除不一定会用的包
	for imp, alias := range imports {
		cs.AddImport(imp, alias)
	}
	if err := cs.AddCode(source); err != nil {
		return cs, err
	}
	for imp, alias := range imports {
		cs.DelImport(imp, alias)
	}
	return cs, nil
}
