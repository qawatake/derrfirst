package dfirst

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const doc = "dfirst is ..."

func NewAnalyzer(pkgPath string, funcName string) *analysis.Analyzer {
	r := runner{
		pkgPath:  pkgPath,
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
	pkgPath  string
	funcName string
	pkgs     map[*types.Package]struct{}
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	r.setPkgs(pass)
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
					if f.Name != r.funcName {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgPath, r.funcName)
						continue
					}
					if _, ok := r.pkgs[pass.TypesInfo.ObjectOf(f).Pkg()]; !ok {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgPath, r.funcName)
						continue
					}
					continue
				case *ast.SelectorExpr:
					if f.Sel.Name != r.funcName {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.packageName(), r.funcName)
						continue
					}
					if _, ok := r.pkgs[pass.TypesInfo.ObjectOf(f.Sel).Pkg()]; !ok {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.packageName(), r.funcName)
						continue
					}
					continue
				}
			default:
				pass.Reportf(decl.Pos(), "should call %s.%s", r.packageName(), r.funcName)
			}
		}
	}
	return nil, nil
}

func (r runner) packageName() string {
	pp := strings.Split(r.pkgPath, "/")
	return pp[len(pp)-1]
}

func (r *runner) setPkgs(pass *analysis.Pass) {
	r.pkgs = make(map[*types.Package]struct{})
	for _, pkg := range pass.Pkg.Imports() {
		if pkg.Path() == r.pkgPath {
			r.pkgs[pkg] = struct{}{}
		}
	}
}
