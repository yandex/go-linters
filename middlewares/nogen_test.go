package middlewares

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/nilness"
)

func TestNoGen(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Nogen(nilness.Analyzer), "nogen")
}
