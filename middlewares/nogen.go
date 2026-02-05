package middlewares

import (
	"golang.org/x/tools/go/analysis"
	"golang.yandex/linters/internal/lintutils"
)

// Nogen adds linting disabling capability for generated files to analyzer
func Nogen(analyzer *analysis.Analyzer) *analysis.Analyzer {
	nogenAnalyzer := *analyzer

	nogenAnalyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
		localPass := *pass

		// swap report func
		localPass.Report = func(d analysis.Diagnostic) {
			if df, ok := lintutils.FileOfReport(&localPass, d); ok && lintutils.IsGenerated(df) {
				return
			}
			pass.Report(d)
		}

		return analyzer.Run(&localPass)
	}

	return &nogenAnalyzer
}
