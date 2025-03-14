package ctxcheck_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"golang.yandex/linters/passes/ctxcheck"
)

func TestCtxArgAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ctxcheck.CtxArgAnalyzer, "a")
}

func TestCtxSaveAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ctxcheck.CtxSaveAnalyzer, "b")
}
