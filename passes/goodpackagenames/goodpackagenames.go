package goodpackagenames

import (
	"go/ast"
	"go/token"
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
				checkImports(pass, genDecl)
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

func checkImports(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		importSpec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}

		importName := importName(importSpec)
		if importName != "" {
			checkImportName(pass, importSpec, importName)
		}
	}
}

func checkImportName(pass *analysis.Pass, importSpec *ast.ImportSpec, importName string) {
	canonicName := canonicImportName(importName)
	if importName != canonicName {
		pass.Reportf(importSpec.Pos(), "invalid import name %s, use %s", importName, canonicName)
	}
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
