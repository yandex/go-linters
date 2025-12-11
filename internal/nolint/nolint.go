package nolint

import (
	"go/ast"
	"go/token"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"

	"golang.yandex/linters/internal/lintutils"
	"golang.yandex/linters/internal/nogen"
)

const (
	Name = "nolint"

	CommentPrefix = "//nolint:"
)

var Analyzer = &analysis.Analyzer{
	Name:             Name,
	Doc:              `removes Nodes under nolint directives for later passes`,
	Run:              run,
	RunDespiteErrors: true,
	ResultType:       reflect.TypeOf(new(Index)),
	Requires: []*analysis.Analyzer{
		nogen.Analyzer,
	},
}

type Index struct {
	idx map[string][]ast.Node
}

// ForLinter returns subset of excluded nodes specifically for given linter
func (i Index) ForLinter(linter string) *LinterIndex {
	li := &LinterIndex{linter: linter}

	if nodes, ok := i.nodesForLinter(linter); ok {
		li.idx = nodes
		sort.Slice(li.idx, func(i, j int) bool {
			return li.idx[i].Pos() < li.idx[j].Pos()
		})
	}

	return li
}

func (i Index) nodesForLinter(linter string) ([]ast.Node, bool) {
	// TODO(buglloc): leave only names in lowercase after migration
	// first try original linter name
	legacyNodes, legacyOK := i.idx[linter]
	// then name in lowercase
	lowerNodes, lowerOK := i.idx[strings.ToLower(linter)]
	return append(legacyNodes, lowerNodes...), legacyOK || lowerOK
}

type LinterIndex struct {
	linter string
	idx    []ast.Node
}

func (l LinterIndex) Excluded(node ast.Node) bool {
	// TODO: binary search here
	for _, n := range l.idx {
		match := false
		switch n.(type) {
		case *ast.File:
			match = node == n
		default:
			match = node.Pos() >= n.Pos() && node.End() <= n.End()
		}

		if match {
			return true
		}
	}
	return false
}

func (l LinterIndex) Contains(pos token.Pos) bool {
	for _, n := range l.idx {
		if n.Pos() <= pos && pos <= n.End() {
			return true
		}
	}
	return false
}

func run(pass *analysis.Pass) (any, error) {
	files := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files).List()

	// gather nolint index
	index := make(map[string][]ast.Node)

	for _, file := range files {
		for _, cg := range file.Comments {
			linters := getNolintNames(cg)
			if len(linters) == 0 {
				continue
			}

			if node, ok := lintutils.CommentNode(cg, file); ok {
				for _, linter := range linters {
					index[linter] = append(index[linter], node)
				}
			}
		}
	}

	return &Index{idx: index}, nil
}

// getNolintNames returns names of linters from `nolint` comment
func getNolintNames(cg *ast.CommentGroup) []string {
	if cg == nil {
		return nil
	}

	var res []string
	for _, cm := range cg.List {
		if strings.HasPrefix(cm.Text, CommentPrefix) && len(cm.Text) > len(CommentPrefix) {
			res = append(res, cm.Text[len(CommentPrefix):])
		}
	}

	return res
}

func CommentForLinter(linter string) string {
	return CommentPrefix + strings.ToLower(linter)
}
