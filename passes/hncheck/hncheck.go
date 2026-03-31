package hncheck

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "hncheck",
	Doc:  `I will not buy this record, it is scratched`,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for ident, obj := range pass.TypesInfo.Defs {
		if obj == nil {
			continue
		}

		prefix := getNotationPrefix(ident.Name, obj.Type().Underlying())
		if prefix != "" {
			pass.Reportf(ident.Pos(), "identifier '%s' contains non-ideomatic notation", ident.Name)
		}
	}

	return nil, nil
}

func getNotationPrefix(name string, typ types.Type) string {
	if typ == nil {
		return ""
	}

	for _, r := range "bfnipsaIT" {
		prefix := string(r)
		if !strings.HasPrefix(name, prefix) || len(name) <= len(prefix) {
			continue
		}

		nextChar := name[len(prefix)]
		if nextChar >= 'A' && nextChar <= 'Z' && matchesType(prefix, typ) {
			return prefix
		}
	}

	return ""
}

// matchesType compares our suspicious prefix against the type-checker's reality.
func matchesType(prefix string, typ types.Type) bool {
	if typ == nil {
		return false
	}

	switch prefix {
	case "b":
		basic, ok := typ.(*types.Basic)
		return ok && basic.Info()&types.IsBoolean != 0
	case "f":
		basic, ok := typ.(*types.Basic)
		return ok && basic.Info()&types.IsFloat != 0
	case "n", "i":
		basic, ok := typ.(*types.Basic)
		return ok && basic.Info()&types.IsInteger != 0
	case "p":
		_, isPtr := typ.(*types.Pointer)
		return isPtr
	case "s":
		basic, ok := typ.(*types.Basic)
		return ok && basic.Info()&types.IsString != 0
	case "a":
		_, isArray := typ.(*types.Array)
		_, isSlice := typ.(*types.Slice)
		return isArray || isSlice
	case "I":
		_, ok := typ.(*types.Interface)
		return ok
	case "T":
		_, ok := typ.(*types.Struct)
		return ok
	}

	return false
}
