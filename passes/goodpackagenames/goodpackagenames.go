package goodpackagenames

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	Name = "goodpackagenames"
	Doc  = `goodpackagenames checks that your packages and imports have correct names.

"Good package names are short and clear. They are lower case, with no under_scores or mixedCaps."

See https://go.dev/blog/package-names for more information.
`

	packageTestSuffix = "_test"
)

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		checkPackageName(pass, file, packageName(file))

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			if genDecl.Tok == token.IMPORT {
				checkImports(pass, file, genDecl)
			}
		}
	}

	return nil, nil
}

func checkPackageName(pass *analysis.Pass, file *ast.File, packageName string) {
	canonicName := canonicPackageName(packageName)
	if packageName != canonicName {
		pass.Reportf(file.Name.End(), "invalid package name %s, use %s", packageName, canonicName)
	}
}

func checkImports(pass *analysis.Pass, file *ast.File, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		importName := importName(importSpec)
		if importName != "" {
			checkImportName(pass, file, importSpec, importName)
		}
	}
}

func checkImportName(pass *analysis.Pass, file *ast.File, importSpec *ast.ImportSpec, importName string) {
	canonicName := canonicImportName(importName)
	if importName == canonicName {
		return
	}

	pass.Report(analysis.Diagnostic{
		Pos:            importSpec.Pos(),
		Message:        fmt.Sprintf("invalid import name %s, use %s", importName, canonicName),
		SuggestedFixes: renameImportFix(pass, file, importSpec, canonicName),
	})
}

// renameImportFix renames the import alias together with all its uses in the file.
// The alias is dropped when the new name matches the imported package name.
// No fix is suggested when the new name is not a valid identifier or is already
// taken by another declaration visible at any of the uses.
func renameImportFix(pass *analysis.Pass, file *ast.File, spec *ast.ImportSpec, newName string) []analysis.SuggestedFix {
	if !token.IsIdentifier(newName) {
		return nil
	}

	pkgName, ok := pass.TypesInfo.Defs[spec.Name].(*types.PkgName)
	if !ok {
		return nil
	}

	var uses []*ast.Ident
	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok && pass.TypesInfo.Uses[ident] == pkgName {
			uses = append(uses, ident)
		}

		return true
	})

	for _, ident := range uses {
		scope := pass.Pkg.Scope().Innermost(ident.Pos())
		if scope == nil {
			return nil
		}

		if _, obj := scope.LookupParent(newName, ident.Pos()); obj != nil {
			return nil
		}
	}

	aliasEdit := analysis.TextEdit{
		Pos:     spec.Name.Pos(),
		End:     spec.Name.End(),
		NewText: []byte(newName),
	}
	if newName == pkgName.Imported().Name() {
		aliasEdit = analysis.TextEdit{
			Pos: spec.Name.Pos(),
			End: spec.Path.Pos(),
		}
	}

	edits := []analysis.TextEdit{aliasEdit}
	for _, ident := range uses {
		edits = append(edits, analysis.TextEdit{
			Pos:     ident.Pos(),
			End:     ident.End(),
			NewText: []byte(newName),
		})
	}

	return []analysis.SuggestedFix{{
		Message:   fmt.Sprintf("Rename import %s to %s", spec.Name.Name, newName),
		TextEdits: edits,
	}}
}

func canonicImportName(name string) string {
	if name == "_" {
		return name
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "")

	return name
}

func canonicPackageName(name string) string {
	testSuffix := ""
	if strings.HasSuffix(name, packageTestSuffix) {
		name = strings.TrimSuffix(name, packageTestSuffix)
		testSuffix = packageTestSuffix
	}

	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "_", "")

	return name + testSuffix
}

func importName(spec *ast.ImportSpec) string {
	if spec.Name == nil {
		return ""
	}
	return spec.Name.Name
}

func packageName(file *ast.File) string {
	return file.Name.Name
}
