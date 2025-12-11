package execinquery_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"golang.yandex/linters/passes/execinquery"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, execinquery.Analyzer, "a")
}
