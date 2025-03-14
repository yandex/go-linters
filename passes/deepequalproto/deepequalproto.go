package deepequalproto

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

type compareFn struct {
	offset          int
	alternativeName string
}

var comparingFn = map[string]compareFn{
	"reflect.DeepEqual": {
		offset:          0,
		alternativeName: "proto.Equal",
	},
	"github.com/stretchr/testify/assert.Equal": {
		offset:          1,
		alternativeName: "assertpb.Equal",
	},
	"github.com/stretchr/testify/assert.Equalf": {
		offset:          1,
		alternativeName: "assertpb.Equalf",
	},
	"github.com/stretchr/testify/require.Equal": {
		offset:          1,
		alternativeName: "requirepb.Equal",
	},
	"github.com/stretchr/testify/require.Equalf": {
		offset:          1,
		alternativeName: "requirepb.Equalf",
	},
	"(*github.com/stretchr/testify/assert.Assertions).Equal": {
		offset:          0,
		alternativeName: "assertpb.Equal",
	},
	"(*github.com/stretchr/testify/assert.Assertions).Equalf": {
		offset:          0,
		alternativeName: "assertpb.Equalf",
	},
	"(*github.com/stretchr/testify/require.Assertions).Equal": {
		offset:          0,
		alternativeName: "requirepb.Equal",
	},
	"(*github.com/stretchr/testify/require.Assertions).Equalf": {
		offset:          0,
		alternativeName: "requirepb.Equalf",
	},
}

var Analyzer = &analysis.Analyzer{
	Name: "deepequalproto",
	Doc:  `deepequalproto checks that protobuf messages are not compared using reflect.DeepEqual`,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

var (
	callFilter = []ast.Node{
		(*ast.CallExpr)(nil),
	}
)

func run(pass *analysis.Pass) (interface{}, error) {
	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	ins.Preorder(callFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
		if !ok {
			return
		}

		compareFn, ok := comparingFn[fn.FullName()]
		if !ok {
			return
		}

		for i := compareFn.offset; i < compareFn.offset+2 && i < len(call.Args); i++ {
			shortName := fn.Pkg().Name() + "." + fn.Name()

			if hasProto(pass, call.Args[i]) {
				pass.ReportRangef(call, "avoid using %s with proto.Message; use %s instead",
					shortName,
					compareFn.alternativeName)
				return
			}
		}
	})

	return nil, nil
}

// hasProto reports whether the type of v contains the proto message.
// See containsProto, below, for the meaning of "contains".
func hasProto(pass *analysis.Pass, v ast.Expr) bool {
	tv, ok := pass.TypesInfo.Types[v]
	if !ok { // no type info, assume good
		return false
	}
	return containsProto(tv.Type)
}

func isProtoType(typ types.Type) bool {
	if t, ok := typ.(*types.Struct); ok {
		for i := 0; i < t.NumFields(); i++ {
			if t.Field(i).Name() == "XXX_unrecognized" {
				return true
			}
		}
	}

	return false
}

func containsProto(typ types.Type) bool {
	// Track types being processed, to avoid infinite recursion.
	// Using types as keys here is OK because we are checking for the identical pointer, not
	// type identity. See analysis/passes/printf/types.go.
	inProgress := make(map[types.Type]bool)

	var check func(t types.Type) bool
	check = func(t types.Type) bool {
		if isProtoType(t) {
			return true
		}

		if inProgress[t] {
			return false
		}
		inProgress[t] = true

		switch t := t.(type) {
		case *types.Pointer:
			return check(t.Elem())
		case *types.Slice:
			return check(t.Elem())
		case *types.Array:
			return check(t.Elem())
		case *types.Map:
			return check(t.Key()) || check(t.Elem())
		case *types.Struct:
			for i := 0; i < t.NumFields(); i++ {
				if check(t.Field(i).Type()) {
					return true
				}
			}
		case *types.Named:
			return check(t.Underlying())
		}
		return false
	}

	return check(typ)
}
