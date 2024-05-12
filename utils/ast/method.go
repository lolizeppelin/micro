package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// GetComment 获取method的注释
func GetComment(method reflect.Method) (string, error) {
	file, _ := runtime.FuncForPC(method.Func.Pointer()).FileLine(method.Func.Pointer())
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	//cname := component.Name()
	name := method.Name
	comment := func() (ret string) {
		f, _ := parser.ParseFile(token.NewFileSet(), file, string(data), parser.ParseComments)
		ast.Inspect(f, func(n ast.Node) bool {
			if astFile, ok := n.(*ast.File); ok {
				for _, decl := range astFile.Decls {
					if astFunc, ok := decl.(*ast.FuncDecl); ok {
						if astFunc.Name.String() == name {
							ret = astFunc.Doc.Text()
							return false
						}
					}
				}
			}
			return true
		})
		return
	}()
	return strings.TrimSpace(comment), nil
}
