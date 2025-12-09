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

func TestForceCasing(t *testing.T) {
	t.Run("no_flag", func(t *testing.T) {
		testdata := analysistest.TestData()
		analysistest.Run(t, testdata, Analyzer, "a")
	})

	t.Run("snake", func(t *testing.T) {
		flagForceCasing = casingSnake
		defer func() {
			flagForceCasing = casingUnknown
		}()

		testdata := analysistest.TestData()
		analysistest.Run(t, testdata, Analyzer, "force_snake")
	})

	t.Run("camel", func(t *testing.T) {
		flagForceCasing = casingCamel
		defer func() {
			flagForceCasing = casingUnknown
		}()

		testdata := analysistest.TestData()
		analysistest.Run(t, testdata, Analyzer, "force_camel")
	})

	t.Run("kebab", func(t *testing.T) {
		flagForceCasing = casingKebab
		defer func() {
			flagForceCasing = casingUnknown
		}()

		testdata := analysistest.TestData()
		analysistest.Run(t, testdata, Analyzer, "force_kebab")
	})
}

func TestDetectCasing(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected stringCasing
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
