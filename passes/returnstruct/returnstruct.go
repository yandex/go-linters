package returnstruct

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"golang.yandex/linters/internal/lintutils"
	"golang.yandex/linters/internal/nogen"
	"golang.yandex/linters/internal/nolint"
)

const (
	Name      = "returnstruct"
	typeError = "error"
)

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc:  Name + ` checks the second half of "Accept Interfaces, Return Structs"`,
	Run:  run,
	Requires: []*analysis.Analyzer{
		nolint.Analyzer,
		nogen.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	nogenFiles := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files)

	nolintIndex := lintutils.ResultOf(pass, nolint.Name).(*nolint.Index)
	nolintNodes := nolintIndex.ForLinter(Name)

	ins := inspector.New(nogenFiles.List())

	// we filter only function declarations
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	ins.Nodes(nodeFilter, func(n ast.Node, push bool) (proceed bool) {
		// do not fall into leaf twice
		if !push {
			return false
		}

		funcDecl := n.(*ast.FuncDecl)

		// skip nolint node
		if nolintNodes.Excluded(funcDecl) {
			return false
		}

		checkFuncDeclSignature(pass, funcDecl)
		return true
	})

	return nil, nil
}

func checkFuncDeclSignature(pass *analysis.Pass, decl *ast.FuncDecl) {
	res := decl.Type.Results
	// function returns no results, skip
	if res == nil || res.NumFields() == 0 {
		return
	}

	for _, param := range res.List {
		typ := pass.TypesInfo.TypeOf(param.Type)

		_, isNamed := typ.(*types.Named)
		if !isNamed || !types.IsInterface(typ) || typ.String() == typeError {
			// we need only named interface types, except built-in `error`
			continue
		}

		pass.Report(analysis.Diagnostic{
			Pos:     param.Pos(),
			Message: fmt.Sprintf("function must return concrete type, not interface %v", typ),
			URL:     "https://medium.com/@cep21/what-accept-interfaces-return-structs-means-in-go-2fe879e25ee8",
		})
	}
}
