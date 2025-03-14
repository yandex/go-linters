package remindercheck_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"golang.yandex/linters/passes/remindercheck"
)

func TestRun(t *testing.T) {
	testdata := analysistest.TestData()

	analysistest.RunWithSuggestedFixes(t, testdata, remindercheck.Analyzer())
}
