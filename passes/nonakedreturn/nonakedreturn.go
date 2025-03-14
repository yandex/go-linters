package nonakedreturn

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `Checks that functions with named results does not have naked returns

See linter tests (testdata/src/a directory) to clarify concrete cases.
`

var Analyzer = &analysis.Analyzer{
	Name:     "nonakedreturn",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// tells which filename we are handling
func getPathToFile(pass *analysis.Pass, file *ast.File) string {
	var result string
	pass.Fset.Iterate(func(f *token.File) bool {
		if int(file.Package) >= f.Base() && int(file.Package) < f.Base()+f.Size() {
			result = f.Name()
			return false
		}
		return true // continue
	})
	return result
}

// handles our hand-made stack of ast-nodes, to determine - which function corresponds given ReturnStmt
func getFuncForReturn(stack []ast.Node, returnStmt *ast.ReturnStmt, fileSet *token.FileSet) ast.Node {
	for i := len(stack) - 1; i != 0; i-- {
		switch stack[i].(type) {
		case *ast.FuncDecl:
			return stack[i]
		case *ast.FuncLit:
			return stack[i]
		}
	}
	panic(fmt.Sprintf("Return statement found with no surrounding function at %s", fileSet.Position(returnStmt.Pos())))
}

func getResultsLen(funcNode ast.Node) int {
	var results *ast.FieldList
	switch funcNode := funcNode.(type) {
	case *ast.FuncLit:
		results = funcNode.Type.Results
	case *ast.FuncDecl:
		results = funcNode.Type.Results
	default:
		panic(fmt.Sprintf("Invalid node type: %T", funcNode))
	}

	if results == nil {
		return 0
	}

	return results.NumFields()
}

type FuncToReturns map[ast.Node][]*ast.ReturnStmt

func extractFuncToReturns(file *ast.File, fileSet *token.FileSet) FuncToReturns {
	funcToReturns := make(FuncToReturns)
	var stack []ast.Node
	ast.Inspect(file, func(node ast.Node) bool {
		// build stack
		if node == nil {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, node)
		}

		// search return
		if returnStmt, ok := node.(*ast.ReturnStmt); ok {
			lastFunc := getFuncForReturn(stack, returnStmt, fileSet)
			if _, ok := funcToReturns[lastFunc]; ok {
				funcToReturns[lastFunc] = append(funcToReturns[lastFunc], returnStmt)
			} else {
				funcToReturns[lastFunc] = make([]*ast.ReturnStmt, 1)
				funcToReturns[lastFunc][0] = returnStmt
			}
		}
		return true
	})
	return funcToReturns
}

func funcName(node ast.Node) string {
	switch node := node.(type) {
	case *ast.FuncLit:
		return "(function literal)"
	case *ast.FuncDecl:
		return node.Name.Name
	default:
		panic(fmt.Sprintf("Invalid node type: %T", node))
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		pathToFile := getPathToFile(pass, file)
		if strings.HasSuffix(pathToFile, "_test.go") || strings.HasSuffix(pathToFile, "_mock.go") {
			continue
		}

		funcToReturns := extractFuncToReturns(file, pass.Fset)
		for currFuncNode, currRets := range funcToReturns {
			resultsNumber := getResultsLen(currFuncNode)
			if resultsNumber == 0 {
				continue
			}

			for i, ret := range currRets {
				if len(ret.Results) == 0 {
					pass.Reportf(
						ret.Pos(),
						"Naked return - %dth return in function %s (should be %d values)",
						i+1,
						funcName(currFuncNode),
						resultsNumber,
					)
				}
			}
		}
	}
	return nil, nil
}
