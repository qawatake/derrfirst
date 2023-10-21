package dfirst

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const doc = "dfirst is ..."

func NewAnalyzer(name string) *analysis.Analyzer {
	pkgName, funcName, err := parse(name)
	if err != nil {
		panic(err)
	}
	r := runner{
		pkgName:  pkgName,
		funcName: funcName,
	}
	return &analysis.Analyzer{
		Name: "dfirst",
		Doc:  doc,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
	}
}

type runner struct {
	pkgName  string
	funcName string
}

func (r runner) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !fn.Name.IsExported() {
				continue
			}
			if fn.Body == nil {
				continue
			}
			if fn.Body.List == nil {
				continue
			}
			first := fn.Body.List[0]
			switch first := first.(type) {
			case *ast.DeferStmt:
				switch f := first.Call.Fun.(type) {
				case *ast.Ident:
					pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName, r.funcName)
				case *ast.SelectorExpr:
					if f.Sel.Name != "Println" || f.X.(*ast.Ident).Name != "fmt" {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName, r.funcName)
					}
				}
			default:
				pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName, r.funcName)
			}
		}
	}
	return nil, nil
}

func parse(name string) (string, string, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 2 {
		return "", "", nil
	}
	splittedPkgPath := strings.Split(parts[0], "/")
	pkgName := splittedPkgPath[len(splittedPkgPath)-1]
	return pkgName, parts[1], nil
}
