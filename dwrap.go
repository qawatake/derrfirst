package dwrap

import (
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/qawatake/dwrap/internal/analysisutil"
	"golang.org/x/tools/go/analysis"
)

const name = "dwrap"
const doc = "dwrap forces every public function to begin with a deferring call of an error wrapping function"
const url = "https://pkg.go.dev/github.com/qawatake/dwrap"

func NewAnalyzer(pkgPath string, funcName string) *analysis.Analyzer {
	r := runner{
		pkgPath:  pkgPath,
		funcName: funcName,
	}
	return &analysis.Analyzer{
		Name: name,
		Doc:  doc,
		URL:  url,
		Run:  r.run,
		Requires: []*analysis.Analyzer{
			commentmap.Analyzer,
		},
	}
}

type runner struct {
	pkgPath  string
	funcName string
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
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
			if !returnsError(pass, fn) {
				continue
			}
			first := fn.Body.List[0]
			switch first := first.(type) {
			case *ast.DeferStmt:
				switch f := first.Call.Fun.(type) {
				case *ast.Ident:
					got := pass.TypesInfo.ObjectOf(f)
					want := analysisutil.ObjectOf(pass, r.pkgPath, r.funcName)
					if got != want {
						pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
						continue
					}
					continue
				case *ast.SelectorExpr:
					got := pass.TypesInfo.ObjectOf(f.Sel)
					want := analysisutil.ObjectOf(pass, r.pkgPath, r.funcName)
					if got != want {
						pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
						continue
					}
					continue
				}
			default:
				pass.Reportf(decl.Pos(), "should call defer %s.%s on the first line", r.pkgName(), r.funcName)
			}
		}
	}
	return nil, nil
}

func (r runner) pkgName() string {
	pp := strings.Split(r.pkgPath, "/")
	return pp[len(pp)-1]
}

func returnsError(pass *analysis.Pass, fn *ast.FuncDecl) bool {
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
