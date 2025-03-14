package goodpackagenames_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"golang.yandex/linters/passes/goodpackagenames"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goodpackagenames.Analyzer, "a/lintexamples")
}
