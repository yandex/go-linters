package lintutils

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// NodeByPos returns node of given position
func NodeByPos(pass *analysis.Pass, pos token.Pos) (node ast.Node, found bool) {
	file, ok := FileOfPos(pass, pos)
	if !ok {
		return nil, false
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil && n.Pos() == pos {
			node = n
			found = true
		}
		return !found
	})

	return
}
