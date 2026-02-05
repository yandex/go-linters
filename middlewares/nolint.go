package middlewares

import (
	"fmt"

	"golang.org/x/tools/go/analysis"
	"golang.yandex/linters/internal/lintutils"
	"golang.yandex/linters/internal/passes/nolint"
)

const (
	nolintDoc = `if you believe this report is false positive, please silence it with %s comment`
)

// Nolint adds linting disabling capability to analyzer
func Nolint(analyzer *analysis.Analyzer) *analysis.Analyzer {
	nolintAnalyzer := *analyzer
	// prepend nolint analyzer to give it maximum priority
	nolintAnalyzer.Requires = append([]*analysis.Analyzer{nolint.Analyzer}, analyzer.Requires...)

	nolintAnalyzer.Run = func(pass *analysis.Pass) (any, error) {
		localPass := *pass

		// gather nolint nodes
		nolintNodes := lintutils.ResultOf(&localPass, nolint.Name).(*nolint.Index).ForLinter(analyzer.Name)

		// swap report func
		localPass.Report = func(d analysis.Diagnostic) {
			// When Analyzer uses ast.Inspect search and reports diagnostics
			// based on (*ast.Node).Pos(), then the reported node could be found
			// and skipping can be determined by relative positions of nodes
			//
			// Actually analyzer could pass any pos in report. If no *ast.Node
			// was found, we could just check if reported position in nolint range
			if dn, found := lintutils.NodeOfReport(&localPass, d); found && nolintNodes.Excluded(dn) {
				return
			} else if !found && nolintNodes.Contains(d.Pos) {
				return
			}

			pass.Report(d)

			d.Message = fmt.Sprintf(nolintDoc, nolint.CommentForLinter(analyzer.Name))
			pass.Report(d)
		}

		return analyzer.Run(&localPass)
	}

	return &nolintAnalyzer
}
