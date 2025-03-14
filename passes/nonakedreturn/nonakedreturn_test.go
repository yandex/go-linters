package nonakedreturn

import (
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

type myTestingWrapper struct {
	original   *testing.T
	errorCases map[string]bool
}

func (m myTestingWrapper) Errorf(format string, args ...any) {
	if format == "%v: unexpected %s: %v" {
		if len(args) > 0 {
			tmp := args[0]
			if tokenPos, ok := tmp.(token.Position); ok {
				if strings.HasPrefix(tokenPos.Filename, "a/error") {
					m.errorCases[tokenPos.Filename] = true
					return
				}
			}
		}
	}
	m.original.Errorf(format, args...)
}

func (m myTestingWrapper) ErrorCount() int {
	return len(m.errorCases)
}

const (
	ExpectedErrorCases = 1
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	testingWrapper := myTestingWrapper{
		original:   t,
		errorCases: make(map[string]bool),
	}
	_ = analysistest.Run(testingWrapper, testdata, Analyzer, "a")
	if testingWrapper.ErrorCount() != ExpectedErrorCases {
		t.Errorf("Got %d error test cases, expected %d", testingWrapper.ErrorCount(), ExpectedErrorCases)
	}
}
