package execinquery

import (
	"go/ast"
	"go/token"
	"maps"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "execinquery is a linter about query string checker in Query function which reads your Go src files and warning it finds"

var (
	commentExp          = regexp.MustCompile(`--[^\n]*\n`)
	multilineCommentExp = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// Analyzer is checking database/sql pkg Query's function
var Analyzer = &analysis.Analyzer{
	Name: "execinquery",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// collect global vars for package
	globalVars := collectGlobalVars(ins)

	// inspect each individual top-level function/method
	funcFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	ins.Preorder(funcFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		inspectFunc(pass, funcDecl, maps.Clone(globalVars))
	})

	return nil, nil
}

func collectGlobalVars(ins *inspector.Inspector) map[string]ast.Expr {
	filter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	vars := make(map[string]ast.Expr)
	ins.Preorder(filter, func(n ast.Node) {
		decl := n.(*ast.GenDecl)
		if decl.Tok != token.VAR && decl.Tok != token.CONST {
			return
		}

		for _, spec := range decl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			for i, ident := range valueSpec.Names {
				if i >= len(valueSpec.Values) {
					break
				}
				vars[ident.Name] = valueSpec.Values[i]
			}
		}
	})

	return vars
}

func inspectFunc(pass *analysis.Pass, funcDecl *ast.FuncDecl, vars map[string]ast.Expr) {
	for n := range ast.Preorder(funcDecl) {
		switch node := n.(type) {
		case *ast.AssignStmt:
			maps.Copy(vars, collectAssignmentVariables(node))
		case *ast.CallExpr:
			inspectCallExpr(pass, node, vars)
		}
	}
}

func collectAssignmentVariables(assn *ast.AssignStmt) map[string]ast.Expr {
	vars := make(map[string]ast.Expr)
	for i, lhs := range assn.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || i >= len(assn.Rhs) {
			continue
		}
		vars[ident.Name] = assn.Rhs[i]
	}

	return vars
}

func inspectCallExpr(pass *analysis.Pass, callExpr *ast.CallExpr, vars map[string]ast.Expr) {
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || !strings.Contains(selector.Sel.Name, "Query") {
		return
	}

	if pass.TypesInfo == nil {
		return
	}

	selectorInfo := pass.TypesInfo.Uses[selector.Sel]
	if selectorInfo == nil || selectorInfo.Pkg() == nil || selectorInfo.Pkg().Path() != "database/sql" {
		return
	}

	replacement := "Exec"
	paramIdx := 0 // the index of the query argument
	if strings.Contains(selector.Sel.Name, "Context") {
		replacement = "ExecContext"
		paramIdx = 1
	}

	if len(callExpr.Args)-1 < paramIdx {
		return
	}

	query := getQueryString(callExpr.Args[paramIdx], vars)
	if query == "" {
		return
	}

	query = strings.TrimSpace(cleanValue(query))
	delim := strings.IndexByte(query, ' ')
	if delim == -1 {
		return
	}

	cmd := query[:delim]
	if strings.EqualFold(cmd, "SELECT") || strings.EqualFold(cmd, "SHOW") {
		return
	}

	// PostgreSQL RETURNING clause makes INSERT/UPDATE/DELETE return rows
	upperQuery := strings.ToUpper(query)
	if strings.Contains(upperQuery, "RETURNING") {
		return
	}

	pass.Reportf(callExpr.Fun.Pos(), "Use %s instead of %s to execute `%s` query", replacement, selector.Sel.Name, cmd)
}

func getQueryString(exp any, vars map[string]ast.Expr) string {
	switch e := exp.(type) {
	case *ast.AssignStmt:
		var b strings.Builder
		for _, stmt := range e.Rhs {
			b.WriteString(cleanValue(getQueryString(stmt, vars)))
		}
		return b.String()

	case *ast.BasicLit:
		return e.Value

	case *ast.ValueSpec:
		var b strings.Builder
		for _, value := range e.Values {
			b.WriteString(cleanValue(getQueryString(value, vars)))
		}
		return b.String()

	case *ast.Ident:
		if assn, ok := vars[e.Name]; ok && assn != nil {
			return getQueryString(assn, vars)
		}
		// Fall back to original declaration
		if e.Obj == nil {
			return ""
		}
		return getQueryString(e.Obj.Decl, vars)

	case *ast.BinaryExpr:
		v := cleanValue(getQueryString(e.X, vars))
		v += cleanValue(getQueryString(e.Y, vars))
		return v
	}

	return ""
}

func cleanValue(s string) string {
	v := strings.NewReplacer(`"`, "", "`", "").Replace(s)
	v = multilineCommentExp.ReplaceAllString(v, "")
	return commentExp.ReplaceAllString(v, "")
}
