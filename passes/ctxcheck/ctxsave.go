package ctxcheck

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var CtxSaveAnalyzer = &analysis.Analyzer{
	Name:     "ctxsave",
	Doc:      `ctxsave ensures the context does not saved as a struct field`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      ctxsave,
}

func ctxsave(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	ins.Preorder(nodeFilter, func(n ast.Node) {
		strct := n.(*ast.StructType)
		if strct.Fields == nil {
			return
		}

		for _, field := range strct.Fields.List {
			typ := types.ExprString(field.Type)
			if typ == "context.Context" {
				pass.Reportf(field.Pos(), "context must not be saved as a struct field")
			}
		}
	})

	return nil, nil
}
