package rewrite

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"

	"github.com/azhai/gozzo/match"
	"golang.org/x/tools/go/ast/astutil"
)

// PosAlt 替换位置
type PosAlt struct {
	Pos, End  token.Position
	Alternate []byte
}

// CodeSource 源码解析器
type CodeSource struct {
	Fileast    *ast.File
	Fileset    *token.FileSet
	Source     []byte
	Alternates []PosAlt // Source 只能替换一次，然后必须重新解析 Fileast
	*printer.Config
}

// NewCodeSource 创造源码解析器
func NewCodeSource() *CodeSource {
	return &CodeSource{
		Fileset: token.NewFileSet(),
		Config: &printer.Config{
			Mode:     printer.TabIndent,
			Tabwidth: 4,
		},
	}
}

// SetSource 替换全部代码，并重新解析
func (cs *CodeSource) SetSource(source []byte) (err error) {
	cs.Source = source
	cs.Fileast, err = parser.ParseFile(cs.Fileset, "", source, parser.ParseComments)
	return
}

// GetContent 获取代码内容
func (cs *CodeSource) GetContent() ([]byte, error) {
	var buf bytes.Buffer
	err := cs.Config.Fprint(&buf, cs.Fileset, cs.Fileast)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// AddCode 增加新代码在原有之后
func (cs *CodeSource) AddCode(code []byte) error {
	content, err := cs.GetContent()
	if err != nil {
		return err
	}
	return cs.SetSource(append(content, code...))
}

// AddStringCode 增加新代码在原有之后
func (cs *CodeSource) AddStringCode(code string) error {
	return cs.AddCode([]byte(code))
}

// GetFirstFileName 如果代码由多个文件组成，返回第一个文件路径
func (cs *CodeSource) GetFirstFileName() string {
	if cs.Fileset == nil || cs.Fileset.Base() <= 1 {
		return ""
	}
	file := cs.Fileset.File(token.Pos(1))
	return file.Name()
}

// GetPackageOffset 获取包名结束位置
func (cs *CodeSource) GetPackageOffset() int {
	if cs.Fileast != nil {
		pos := cs.Fileast.Name.End()
		return cs.Fileset.PositionFor(pos, false).Offset
	}
	return 0
}

// GetPackage 获取包名
func (cs *CodeSource) GetPackage() string {
	if cs.Fileast == nil {
		return ""
	}
	return cs.Fileast.Name.Name
}

// SetPackage 设置新的包名
func (cs *CodeSource) SetPackage(name string) (err error) {
	if cs.Fileast == nil {
		code := fmt.Sprintf("package %s", name)
		err = cs.SetSource([]byte(code))
	} else {
		cs.Fileast.Name.Name = name
	}
	return
}

// AddImport 增加一个import
func (cs *CodeSource) AddImport(path, alias string) bool {
	return astutil.AddNamedImport(cs.Fileset, cs.Fileast, alias, path)
}

// DelImport 删除一个import
func (cs *CodeSource) DelImport(path, alias string) bool {
	if astutil.UsesImport(cs.Fileast, path) {
		return false
	}
	return astutil.DeleteNamedImport(cs.Fileset, cs.Fileast, alias, path)
}

// CleanImports 整理全部import代码
func (cs *CodeSource) CleanImports() (removes int) {
	for _, groups := range astutil.Imports(cs.Fileset, cs.Fileast) {
		for _, imp := range groups {
			var path, alias string
			if imp.Name != nil {
				alias = imp.Name.Name
			}
			path = strings.Trim(imp.Path.Value, "\"")
			if alias != "" || !strings.Contains(path, "/v") {
				continue
			}
			idx := strings.LastIndex(path, "/v")
			if match.IsDigit(path[idx+2:]) { // astutil对带版本号的path判断失误
				// fmt.Println(path, ":", alias)
				// fmt.Println(path, " -> ", astutil.UsesImport(cs.Fileast, path))
				// fmt.Println(path[:idx], " -> ", astutil.UsesImport(cs.Fileast, path[:idx]))
				// fmt.Println("===============================")
				continue
			}
			if cs.DelImport(path, alias) {
				removes++
			}
		}
	}
	return
}

// GetNodeCode 获得节点代码内容
func (cs *CodeSource) GetNodeCode(node ast.Node) string {
	// 请先保证 node 不是 nil
	pos := cs.Fileset.PositionFor(node.Pos(), false)
	end := cs.Fileset.PositionFor(node.End(), false)
	return string(cs.Source[pos.Offset:end.Offset])
}

// GetFieldCode 获得类成员代码内容
func (cs *CodeSource) GetFieldCode(node *DeclNode, i int) string {
	if i < 0 {
		i += len(node.Fields)
	}
	ffcode := cs.GetNodeCode(node.Fields[i])
	return strings.TrimSpace(ffcode)
}

// GetComment 获取注释
func (cs *CodeSource) GetComment(c *ast.CommentGroup, trim bool) string {
	if c == nil {
		return ""
	}
	comment := cs.GetNodeCode(c)
	if trim {
		comment = TrimComment(comment)
	}
	return comment
}

// AddReplace 将两个节点以及中间的部分，使用新内容代替
func (cs *CodeSource) AddReplace(first, last ast.Node, code string) {
	// 请先保证 first, last 不是 nil
	pos := cs.Fileset.PositionFor(first.Pos(), false)
	end := cs.Fileset.PositionFor(last.End(), false)
	alt := PosAlt{Pos: pos, End: end, Alternate: []byte(code)}
	cs.Alternates = append(cs.Alternates, alt)
}

// AltSource 改写源码，应用事先准备的可代替代码Alternates
func (cs *CodeSource) AltSource() ([]byte, bool) {
	if len(cs.Alternates) == 0 {
		return cs.Source, false
	}
	sort.Slice(cs.Alternates, func(i, j int) bool {
		return cs.Alternates[i].Pos.Offset < cs.Alternates[j].Pos.Offset
	})
	var chunks [][]byte
	start, stop := 0, 0
	for _, alt := range cs.Alternates {
		start = alt.Pos.Offset
		chunks = append(chunks, cs.Source[stop:start])
		chunks = append(chunks, alt.Alternate)
		stop = alt.End.Offset
	}
	if stop < len(cs.Source) {
		chunks = append(chunks, cs.Source[stop:])
	}
	cs.Alternates = make([]PosAlt, 0)
	return bytes.Join(chunks, nil), true
}

// WriteTo 美化代码并保存到文件
func (cs *CodeSource) WriteTo(filename string) error {
	code, err := cs.GetContent()
	if err != nil {
		return err
	}
	_, err = WriteGolangFilePrettify(filename, code)
	return err
}

// ResetImports 重新注入声明，并美化代码
func (cs *CodeSource) ResetImports(filename string, imports map[string]string) error {
	if code, chg := cs.AltSource(); chg {
		_ = cs.SetSource(code)
	}
	pkg, offset := cs.GetPackage(), cs.GetPackageOffset()
	source, err := FormatGolangCode(cs.Source[offset:])
	if err != nil {
		return err
	}
	var obj *CodeSource
	obj, err = WriteWithImports(pkg, source, imports)
	if err != nil {
		return err
	}
	if filename == "" {
		filename = cs.GetFirstFileName()
	}
	return obj.WriteTo(filename)
}
