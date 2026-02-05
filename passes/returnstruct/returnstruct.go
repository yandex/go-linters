package returnstruct

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	Name      = "returnstruct"
	typeError = "error"
)

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc:  Name + ` checks the second half of "Accept Interfaces, Return Structs"`,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	ins := inspector.New(pass.Files)

	// we filter only function declarations
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	ins.Nodes(nodeFilter, func(n ast.Node, push bool) (proceed bool) {
		// do not fall into leaf twice
		if !push {
			return false
		}

		checkFuncDeclSignature(pass, n.(*ast.FuncDecl))
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
