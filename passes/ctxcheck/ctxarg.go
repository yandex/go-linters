package ctxcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var CtxArgAnalyzer = &analysis.Analyzer{
	Name:     "ctxarg",
	Doc:      `ctxarg ensures the context parameter is always the first received argument`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      ctxarg,
}

func ctxarg(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncType)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		function := n.(*ast.FuncType)

		for key, f := range function.Params.List {
			typ := pass.TypesInfo.TypeOf(f.Type)
			if typ.String() == `context.Context` && key > 0 {
				pass.Reportf(function.Pos(), "context parameter must be supplied as first argument of function")
			}
		}
	})

	return nil, nil
}
