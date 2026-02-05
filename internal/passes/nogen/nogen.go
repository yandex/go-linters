package nogen

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/analysis"

	"golang.yandex/linters/internal/lintutils"
)

const (
	Name = "nogen"
)

var Analyzer = &analysis.Analyzer{
	Name:             Name,
	Doc:              `remove generated files for later passes`,
	Run:              run,
	RunDespiteErrors: true,
	ResultType:       reflect.TypeFor[*Files](),
}

type Files struct {
	list      []*ast.File
	generated []*ast.File
}

func (f *Files) List() []*ast.File {
	return f.list
}

func (f *Files) Generated() []*ast.File {
	return f.generated
}

func run(pass *analysis.Pass) (any, error) {
	nonGenFiles := make([]*ast.File, 0, len(pass.Files)/2)
	genFiles := make([]*ast.File, 0, len(pass.Files)/2)

	for _, file := range pass.Files {
		if !lintutils.IsGenerated(file) {
			nonGenFiles = append(nonGenFiles, file)
		} else {
			genFiles = append(genFiles, file)
		}
	}

	return &Files{list: nonGenFiles, generated: genFiles}, nil
}
