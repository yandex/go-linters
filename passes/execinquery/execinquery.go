package execinquery

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "execinquery is a linter about query string checker in Query function which reads your Go src files and warning it finds"

// Analyzer is checking database/sql pkg Query's function
var Analyzer = &analysis.Analyzer{
	Name: "execinquery",
	Doc:  doc,
	Run:  newLinter().run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

type linter struct {
	commentExp          *regexp.Regexp
	multilineCommentExp *regexp.Regexp
	// varAssignments tracks the most recent assignment for each variable
	// key: variable name, value: map of position to assigned value
	varAssignments map[string]map[token.Pos]ast.Expr
}

func newLinter() *linter {
	return &linter{
		commentExp:          regexp.MustCompile(`--[^\n]*\n`),
		multilineCommentExp: regexp.MustCompile(`(?s)/\*.*?\*/`),
		varAssignments:      make(map[string]map[token.Pos]ast.Expr),
	}
}

func (l *linter) run(pass *analysis.Pass) (any, error) {
	result := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// First pass: collect all variable assignments
	assignFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	result.Preorder(assignFilter, func(n ast.Node) {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return
		}

		for i, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || i >= len(assign.Rhs) {
				continue
			}

			if l.varAssignments[ident.Name] == nil {
				l.varAssignments[ident.Name] = make(map[token.Pos]ast.Expr)
			}
			l.varAssignments[ident.Name][assign.Pos()] = assign.Rhs[i]
		}
	})

	// Second pass: check Query calls
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	result.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CallExpr:
			selector, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}

			if pass.TypesInfo == nil || pass.TypesInfo.Uses[selector.Sel] == nil || pass.TypesInfo.Uses[selector.Sel].Pkg() == nil {
				return
			}

			if pass.TypesInfo.Uses[selector.Sel].Pkg().Path() != "database/sql" {
				return
			}

			if !strings.Contains(selector.Sel.Name, "Query") {
				return
			}

			replacement := "Exec"
			var i int // the index of the query argument
			if strings.Contains(selector.Sel.Name, "Context") {
				replacement += "Context"
				i = 1
			}

			if len(n.Args) <= i {
				return
			}

			query := l.getQueryString(n.Args[i], n.Pos())
			if query == "" {
				return
			}

			query = strings.TrimSpace(l.cleanValue(query))
			parts := strings.SplitN(query, " ", 2)
			cmd := strings.ToUpper(parts[0])

			if strings.HasPrefix(cmd, "SELECT") || strings.HasPrefix(cmd, "SHOW") {
				return
			}

			// PostgreSQL RETURNING clause makes INSERT/UPDATE/DELETE return rows
			upperQuery := strings.ToUpper(query)
			if strings.Contains(upperQuery, "RETURNING") {
				return
			}

			pass.Reportf(n.Fun.Pos(), "Use %s instead of %s to execute `%s` query", replacement, selector.Sel.Name, cmd)
		}
	})

	return nil, nil
}

func (l linter) cleanValue(s string) string {
	v := strings.NewReplacer(`"`, "", "`", "").Replace(s)

	v = l.multilineCommentExp.ReplaceAllString(v, "")

	return l.commentExp.ReplaceAllString(v, "")
}

func (l *linter) getQueryString(exp any, callPos token.Pos) string {
	switch e := exp.(type) {
	case *ast.AssignStmt:
		var v string
		for _, stmt := range e.Rhs {
			v += l.cleanValue(l.getQueryString(stmt, callPos))
		}
		return v

	case *ast.BasicLit:
		return e.Value

	case *ast.ValueSpec:
		var v string
		for _, value := range e.Values {
			v += l.cleanValue(l.getQueryString(value, callPos))
		}
		return v

	case *ast.Ident:
		// Check for reassignments first
		if assignments, ok := l.varAssignments[e.Name]; ok {
			// Find the most recent assignment before the call position
			var latestPos token.Pos
			var latestExpr ast.Expr
			for pos, expr := range assignments {
				if pos < callPos && pos > latestPos {
					latestPos = pos
					latestExpr = expr
				}
			}
			if latestExpr != nil {
				return l.getQueryString(latestExpr, callPos)
			}
		}

		// Fall back to original declaration
		if e.Obj == nil {
			return ""
		}
		return l.getQueryString(e.Obj.Decl, callPos)

	case *ast.BinaryExpr:
		v := l.cleanValue(l.getQueryString(e.X, callPos))
		v += l.cleanValue(l.getQueryString(e.Y, callPos))
		return v
	}

	return ""
}
