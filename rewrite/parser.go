package rewrite

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

	"github.com/azhai/gozzo/match"
	"github.com/pkg/errors"
)

// GetNameList 获取多个标识的名称列表
func GetNameList(ids []*ast.Ident) (names []string) {
	for _, id := range ids {
		names = append(names, id.Name)
	}
	return
}

// FieldNode 类成员节点
type FieldNode struct {
	Names   []string
	Comment *ast.CommentGroup
	*ast.Field
}

// GetTag 类成员的标签
func (n *FieldNode) GetTag() reflect.StructTag {
	tag := ""
	if n.Tag != nil {
		tag = strings.Trim(n.Tag.Value, "` ")
	}
	return reflect.StructTag(tag)
}

// DeclNode 声明节点
type DeclNode struct {
	Token   token.Token
	Kinds   []string
	Names   []string
	Fields  []*FieldNode
	Comment *ast.CommentGroup
	Offset  int
	ast.Decl
}

// NewDeclNode 创建声明节点
func NewDeclNode(decl ast.Decl, offset int, position token.Pos) (n *DeclNode, err error) {
	n = &DeclNode{Decl: decl, Offset: offset}
	if gen, ok := decl.(*ast.GenDecl); ok {
		n.Comment = gen.Doc
		n.Token = gen.Tok // IMPORT, CONST, TYPE, VAR
		n.Kinds = []string{n.Token.String()}
		gen.TokPos += position
		err = n.ParseGenDecl(gen)
	} else if fun, ok := decl.(*ast.FuncDecl); ok {
		n.Comment = fun.Doc
		n.Token = token.FUNC
		n.Kinds = []string{n.Token.String()}
		fun.Type.Func += position
		err = n.ParseFunDecl(fun)
	} else {
		n.Token = token.ILLEGAL
	}
	return
}

// GetKind 声明节点种类
func (n *DeclNode) GetKind() string {
	return strings.Join(n.Kinds, ".")
}

// GetName 声明节点名称
func (n *DeclNode) GetName() string {
	return strings.Join(n.Names, ", ")
}

// ParseFunDecl //解析函数节点
func (n *DeclNode) ParseFunDecl(fun *ast.FuncDecl) (err error) {
	n.Names = []string{fun.Name.Name}
	return
}

// ParseGenDecl 解析类声明节点
func (n *DeclNode) ParseGenDecl(gen *ast.GenDecl) (err error) {
	for _, spec := range gen.Specs {
		if s, ok := spec.(*ast.ValueSpec); ok {
			n.Names = GetNameList(s.Names)
		} else if s, ok := spec.(*ast.ImportSpec); ok {
			if s.Name != nil {
				n.Names = []string{s.Name.Name}
			}
		} else if s, ok := spec.(*ast.TypeSpec); ok {
			if s.Name != nil {
				n.Names = []string{s.Name.Name}
			}
			if t, ok := s.Type.(*ast.StructType); ok {
				n.Kinds = append(n.Kinds, "struct")
				err = n.ParseStruct(t.Fields.List)
			}
		}
	}
	return
}

// ParseStruct 解析类成员节点
func (n *DeclNode) ParseStruct(fields []*ast.Field) (err error) {
	for _, f := range fields {
		n.Fields = append(n.Fields, &FieldNode{
			Names:   GetNameList(f.Names),
			Comment: f.Comment,
			Field:   f,
		})
	}
	return
}

// CodeParser 代码节点解析器
type CodeParser struct {
	DeclNodes   []*DeclNode
	DeclIndexes map[string][]int
	*CodeSource
}

// NewCodeParser 创建代码节点解析器
func NewCodeParser() *CodeParser {
	return &CodeParser{
		DeclIndexes: make(map[string][]int),
		CodeSource:  NewCodeSource(),
	}
}

// NewFileParser 从文件创建解析器
func NewFileParser(filename string) (cp *CodeParser, err error) {
	cp = NewCodeParser()
	if cp.Source, err = os.ReadFile(filename); err != nil {
		err = errors.Wrap(err, "Read file failed "+filename)
		return
	}
	cp.Fileast, err = parser.ParseFile(cp.Fileset, filename, nil, parser.ParseComments)
	if err != nil {
		err = errors.Wrap(err, "Parse file failed "+filename)
	}
	return
}

// NewSourceParser 从代码创建解析器
func NewSourceParser(source []byte) (cp *CodeParser, err error) {
	cp = NewCodeParser()
	err = cp.SetSource(source)
	return
}

// ParseDecls 解析指定种类的声明
func (cp *CodeParser) ParseDecls(kind string, limit int) bool {
	offset := len(cp.DeclNodes)
	for i, decl := range cp.Fileast.Decls[offset:] {
		index := i + offset
		node, _ := NewDeclNode(decl, index, 0)
		cp.DeclNodes = append(cp.DeclNodes, node)
		tokname := node.Token.String()
		cp.DeclIndexes[tokname] = append(cp.DeclIndexes[tokname], index)
		cp.DeclIndexes[""] = append(cp.DeclIndexes[""], index)
		if limit < 0 {
			continue
		}
		if kind == "" || kind == tokname {
			if len(cp.DeclIndexes[kind]) >= limit {
				return true
			}
		}
	}
	return false
}

// AllDeclNode 获取指定种类的所有节点
func (cp *CodeParser) AllDeclNode(kind string) []*DeclNode {
	var nodes []*DeclNode
	cp.ParseDecls(kind, -1)
	if idxes, ok := cp.DeclIndexes[kind]; ok {
		for _, idx := range idxes {
			nodes = append(nodes, cp.DeclNodes[idx])
		}
	}
	return nodes
}

// FindDeclNode 根据名称规则查找
func (cp *CodeParser) FindDeclNode(kind string, wildcards ...string) []*DeclNode {
	var nodes []*DeclNode
	matchers := match.NewGlobs(wildcards)
	for _, node := range cp.AllDeclNode(kind) {
		if matchers.MatchAny(node.GetName(), true) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetDeclNode 获取指定种类的一个节点
func (cp *CodeParser) GetDeclNode(kind string, offset int) *DeclNode {
	count := offset + 1
	if idxes, ok := cp.DeclIndexes[kind]; ok {
		if offset >= 0 && len(idxes) > offset {
			index := idxes[offset]
			return cp.DeclNodes[index]
		}
		count -= len(idxes)
	}
	success := cp.ParseDecls(kind, count)
	if idxes, ok := cp.DeclIndexes[kind]; ok {
		if success {
			index := idxes[offset]
			return cp.DeclNodes[index]
		}
		if offset < 0 && offset >= 0-len(idxes) {
			index := idxes[offset+len(idxes)]
			return cp.DeclNodes[index]
		}
	}
	return nil
}
