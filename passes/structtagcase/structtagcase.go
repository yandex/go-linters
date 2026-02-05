package structtagcase

import (
	"flag"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	flags.Var(&flagForceCasing, "force-casing", "force specific case to be used in struct tags: snake, camel, kebab")
}

const Name = "structtagcase"

type stringCasing string

const (
	casingUnknown stringCasing = ""
	casingSnake   stringCasing = "snake"
	casingCamel   stringCasing = "camel"
	casingKebab   stringCasing = "kebab"
	casingMixed   stringCasing = "mixed"
)

func (s *stringCasing) Set(v string) error {
	switch stringCasing(v) {
	case casingSnake, casingCamel, casingKebab:
		*s = stringCasing(v)
	}
	return nil
}

func (s stringCasing) String() string {
	switch s {
	case casingSnake, casingCamel, casingKebab, casingMixed:
		return string(s)
	default:
		return "unknown"
	}
}

var (
	knownKeys = []string{"json", "bson", "xml", "yaml"}

	flags           flag.FlagSet
	flagForceCasing stringCasing
)

var Analyzer = &analysis.Analyzer{
	Name:  Name,
	Doc:   `structtagcase checks that you use consistent name case in struct tags`,
	Run:   run,
	Flags: flags,
}

func run(pass *analysis.Pass) (any, error) {
	ins := inspector.New(pass.Files)

	// filter only function calls.
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	ins.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		// do not fall into leaf twice
		if !push {
			return false
		}

		checkTagsCasing(pass, n.(*ast.StructType), flagForceCasing)
		return true
	})

	return nil, nil
}

func checkTagsCasing(pass *analysis.Pass, node *ast.StructType, forcedCase stringCasing) {
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

			tagCasing := detectCasing(name)
			if forcedCase != casingUnknown && tagCasing != casingUnknown && tagCasing != forcedCase {
				pass.Reportf(field.End(), "%s struct tag must be in %s case: %s", tagKey, forcedCase, name)
				continue
			}

			if tagCasing == casingMixed {
				pass.Reportf(field.End(), "unknown casing in %s struct tag: %s", tagKey, name)
				continue
			}

			// store first detected casing unconditionally
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

func detectCasing(value string) stringCasing {
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
