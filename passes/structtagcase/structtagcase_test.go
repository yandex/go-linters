package structtagcase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}

func TestDetectCasing(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected int
	}{
		{"snake_case", "ololo_trololo", casingSnake},
		{"camel_case", "ololoTrololo", casingCamel},
		{"unknown_case", "ololo", casingUnknown},
		{"kebab_case", "ololo-trololo", casingKebab},
		{"some_strange_thing", "ololo_trololoShimbaBoomba", casingMixed},
		{"some_strange_thing", "ololo-trololoShimbaBoomba", casingMixed},
		{"some_strange_thing", "ololo-trololo-shimbaBoomba", casingMixed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			casing := detectCasing(tc.value)
			assert.Equalf(t, tc.expected, casing, "value: %s", tc.value)
		})
	}
}
