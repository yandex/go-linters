package lintutils

import (
	"go/ast"
	"strings"
)

// HasCommentPrefix checks if Comment group has particular prefix in any comment line
func HasCommentPrefix(cg *ast.CommentGroup, prefix string) bool {
	if cg == nil {
		return false
	}

	for _, cm := range cg.List {
		if strings.HasPrefix(cm.Text, prefix) && len(cm.Text) > len(prefix) {
			return true
		}
	}

	return false
}

// HasComment checks if Comment group has particular comment line
func HasComment(cg *ast.CommentGroup, comment string) bool {
	if cg == nil {
		return false
	}

	for _, cm := range cg.List {
		if cm.Text == comment {
			return true
		}
	}

	return false
}

// CommentNode returns next node after given comment
func CommentNode(cg *ast.CommentGroup, file *ast.File) (node ast.Node, found bool) {
	if cg == nil || file == nil {
		return
	}

	if cg.Pos() < file.FileStart || cg.End() > file.End() {
		return
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if !found && n != nil && n.Pos() > cg.Pos() {
			// It is possible that multiple nodes have same .Pos() after cg.
			// In that situation the more broad one is the result, for example:
			//
			// // my comment
			// res := Func()
			// ^
			// It is a .Pos() of both *ast.Ident (res itself) and
			// *ast.AssignStmt, but it assumed that such comment
			// is related to whole assignment, not just the ident
			//
			// ast.Inspect walks by depth-first order, that's why first
			// found node is the result
			node = n
			found = true
		}
		return !found
	})

	return
}

// NodeComments returns node comments
func NodeComments(node ast.Node, file *ast.File) (cg *ast.CommentGroup, found bool) {
	if node == nil || file == nil {
		return
	}

	if node.Pos() < file.Pos() || node.End() > file.End() {
		return
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil && n.Pos() < node.Pos() {
			node = n
			found = true
		}
		return !found
	})

	return
}
