package middlewares

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/nilness"
)

func TestNoLint(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Nolint(nilness.Analyzer), "nolint")
}
