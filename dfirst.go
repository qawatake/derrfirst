package dfirst

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
)

const doc = "dfirst requires that every public function begins by deferring a call to a specific function"

const name = "dfirst"

func NewAnalyzer(pkgPath string, funcName string) *analysis.Analyzer {
	r := runner{
		pkgPath:  pkgPath,
		funcName: funcName,
	}
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			commentmap.Analyzer,
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
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if cmaps.IgnorePos(decl.Pos(), name) {
				continue
			}
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
			if !returnError(pass, fn) {
				continue
			}
			first := fn.Body.List[0]
			switch first := first.(type) {
			case *ast.DeferStmt:
				switch f := first.Call.Fun.(type) {
				case *ast.Ident:
					if f.Name != r.funcName {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName(), r.funcName)
						continue
					}
					if _, ok := r.pkgs[pass.TypesInfo.ObjectOf(f).Pkg()]; !ok {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName(), r.funcName)
						continue
					}
					continue
				case *ast.SelectorExpr:
					if f.Sel.Name != r.funcName {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName(), r.funcName)
						continue
					}
					if _, ok := r.pkgs[pass.TypesInfo.ObjectOf(f.Sel).Pkg()]; !ok {
						pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName(), r.funcName)
						continue
					}
					continue
				}
			default:
				pass.Reportf(decl.Pos(), "should call %s.%s", r.pkgName(), r.funcName)
			}
		}
	}
	return nil, nil
}

func (r runner) pkgName() string {
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

func returnError(pass *analysis.Pass, fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, arg := range fn.Type.Results.List {
		if types.Identical(pass.TypesInfo.TypeOf(arg.Type), types.Universe.Lookup("error").Type()) {
			return true
		}
	}
	return false
}
