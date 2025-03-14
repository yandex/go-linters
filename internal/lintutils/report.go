package lintutils

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// FileOfReport returns file in which report occurs
func FileOfReport(pass *analysis.Pass, d analysis.Diagnostic) (file *ast.File, found bool) {
	return FileOfPos(pass, d.Pos)
}

// NodeOfReport finds node of report
func NodeOfReport(pass *analysis.Pass, d analysis.Diagnostic) (node ast.Node, found bool) {
	return NodeByPos(pass, d.Pos)
}
