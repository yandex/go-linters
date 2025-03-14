package structtagcase

import (
	"go/ast"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"golang.yandex/linters/internal/lintutils"
	"golang.yandex/linters/internal/nogen"
	"golang.yandex/linters/internal/nolint"
)

const (
	Name = "structtagcase"

	casingUnknown = iota
	casingSnake
	casingCamel
	casingKebab
	casingMixed
)

var (
	knownKeys = []string{"json", "bson", "xml", "yaml"}
)

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc:  `structtagcase checks that you use consistent name case in struct tags`,
	Run:  run,
	Requires: []*analysis.Analyzer{
		nolint.Analyzer,
		nogen.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	nogenFiles := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files)

	nolintIndex := lintutils.ResultOf(pass, nolint.Name).(*nolint.Index)
	nolintNodes := nolintIndex.ForLinter(Name)

	ins := inspector.New(nogenFiles.List())

	// filter only function calls.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	ins.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		// do not fall into leaf twice
		if !push {
			return false
		}

		structNode := n.(*ast.StructType)

		// skip nolint node
		if nolintNodes.Excluded(structNode) {
			return false
		}

		checkTagsCasing(pass, structNode)

		return true
	})

	return nil, nil
}

func checkTagsCasing(pass *analysis.Pass, node *ast.StructType) {
	for _, tagKey := range knownKeys {
		// start casing for struct key
		keyCasing := casingUnknown

		for _, field := range node.Fields.List {
			if field.Tag == nil {
				continue
			}

			rawTag, _ := strconv.Unquote(field.Tag.Value)
			if rawTag == "" {
				continue
			}

			structTag, ok := reflect.StructTag(rawTag).Lookup(tagKey)
			if !ok {
				continue
			}

			name := extractTagName(structTag)
			if name == "" || name == "-" {
				continue
			}

			// store first detected casing unconditionally
			tagCasing := detectCasing(name)
			if tagCasing == casingMixed {
				pass.Reportf(field.End(), "unknown casing in %s struct tag: %s", tagKey, name)
				break
			}

			if keyCasing == casingUnknown {
				keyCasing = tagCasing
				continue
			}

			if tagCasing != casingUnknown && tagCasing != keyCasing {
				pass.Reportf(field.End(), "inconsistent text case in %s struct tag: %s", tagKey, name)
			}
		}
	}
}

func extractTagName(value string) string {
	name := value
	idx := strings.Index(value, ",")
	if idx != -1 {
		name = value[:idx]
	}
	return name
}

func detectCasing(value string) int {
	var hasUnderscore, hasDash, hasLowercase, hasUppercase bool

	for _, r := range value {
		// we have all we need - stop here
		if hasUnderscore && hasDash && hasLowercase && hasUppercase {
			break
		}

		if r == '_' {
			hasUnderscore = true
			continue
		}
		if r == '-' {
			hasDash = true
			continue
		}
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				hasUppercase = true
				continue
			}
			if unicode.IsLower(r) {
				hasLowercase = true
				continue
			}
		}
	}

	// mixed case
	if hasLowercase && hasUppercase && (hasUnderscore || hasDash) {
		return casingMixed
	}
	// snake
	if hasUnderscore && (hasLowercase || hasUppercase) {
		return casingSnake
	}
	// kebab
	if hasDash && (hasLowercase || hasUppercase) {
		return casingKebab
	}
	// camel
	if !hasUnderscore && hasLowercase && hasUppercase {
		return casingCamel
	}
	// single word probably
	return casingUnknown
}
