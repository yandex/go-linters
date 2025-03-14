package copyproto

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:      "copyproto",
	Doc:       `copyproto checks that protobuf messages are not copied`,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&IsGoGoPkg{}},
	Run:       run,
}

type IsGoGoPkg struct{}

func (*IsGoGoPkg) AFact() {}

func (*IsGoGoPkg) String() string {
	return "isgogo"
}

func format(fset *token.FileSet, x ast.Expr) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, x)
	return b.String()
}

func markGoGoPkg(pass *analysis.Pass) {
	for _, f := range pass.Files {
		if len(f.Comments) == 0 {
			continue
		}

		for _, comment := range f.Comments[0].List {
			if strings.Contains(comment.Text, "protoc-gen-gogo") {
				pass.ExportPackageFact(&IsGoGoPkg{})
				return
			}
		}
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	markGoGoPkg(pass)

	ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
		(*ast.CallExpr)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.GenDecl)(nil),
		(*ast.RangeStmt)(nil),
		(*ast.ReturnStmt)(nil),
	}
	ins.Preorder(nodeFilter, func(node ast.Node) {
		switch node := node.(type) {
		case *ast.RangeStmt:
			checkCopyProtoRange(pass, node)
		case *ast.FuncDecl:
			checkCopyProtoFunc(pass, node.Name.Name, node.Recv, node.Type)
		case *ast.FuncLit:
			checkCopyProtoFunc(pass, "func", nil, node.Type)
		case *ast.CallExpr:
			checkCopyProtoCallExpr(pass, node)
		case *ast.AssignStmt:
			checkCopyProtoAssign(pass, node)
		case *ast.GenDecl:
			checkCopyProtoGenDecl(pass, node)
		case *ast.CompositeLit:
			checkCopyProtoCompositeLit(pass, node)
		case *ast.ReturnStmt:
			checkCopyProtoReturnStmt(pass, node)
		}
	})
	return nil, nil
}

// checkCopyProtoAssign checks whether an assignment
// copies a proto.
func checkCopyProtoAssign(pass *analysis.Pass, as *ast.AssignStmt) {
	for i, x := range as.Rhs {
		if path := protoPathRhs(pass, x); path != nil {
			pass.ReportRangef(x, "assignment copies proto value to %v: %v", format(pass.Fset, as.Lhs[i]), path)
		}
	}
}

// checkCopyProtoGenDecl checks whether proto is copied
// in variable declaration.
func checkCopyProtoGenDecl(pass *analysis.Pass, gd *ast.GenDecl) {
	if gd.Tok != token.VAR {
		return
	}
	for _, spec := range gd.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		for i, x := range valueSpec.Values {
			if path := protoPathRhs(pass, x); path != nil {
				pass.ReportRangef(x, "variable declaration copies proto value to %v: %v", valueSpec.Names[i].Name, path)
			}
		}
	}
}

// checkCopyProtoCompositeLit detects proto copy inside a composite literal
func checkCopyProtoCompositeLit(pass *analysis.Pass, cl *ast.CompositeLit) {
	for _, x := range cl.Elts {
		if node, ok := x.(*ast.KeyValueExpr); ok {
			x = node.Value
		}
		if path := protoPathRhs(pass, x); path != nil {
			pass.ReportRangef(x, "literal copies proto value from %v: %v", format(pass.Fset, x), path)
		}
	}
}

// checkCopyProtoReturnStmt detects proto copy in return statement
func checkCopyProtoReturnStmt(pass *analysis.Pass, rs *ast.ReturnStmt) {
	for _, x := range rs.Results {
		if path := protoPathRhs(pass, x); path != nil {
			pass.ReportRangef(x, "return copies proto value: %v", path)
		}
	}
}

// checkCopyProtoCallExpr detects proto copy in the arguments to a function call
func checkCopyProtoCallExpr(pass *analysis.Pass, ce *ast.CallExpr) {
	var id *ast.Ident
	switch fun := ce.Fun.(type) {
	case *ast.Ident:
		id = fun
	case *ast.SelectorExpr:
		id = fun.Sel
	}
	if fun, ok := pass.TypesInfo.Uses[id].(*types.Builtin); ok {
		switch fun.Name() {
		case "new", "len", "cap", "Sizeof":
			return
		}
	}
	for _, x := range ce.Args {
		if path := protoPathRhs(pass, x); path != nil {
			pass.ReportRangef(x, "call of %s copies proto value: %v", format(pass.Fset, ce.Fun), path)
		}
	}
}

// checkCopyProtoFunc checks whether a function might
// inadvertently copy a proto, by checking whether
// its receiver, parameters, or return values
// are protos.
func checkCopyProtoFunc(pass *analysis.Pass, name string, recv *ast.FieldList, typ *ast.FuncType) {
	if recv != nil && len(recv.List) > 0 {
		expr := recv.List[0].Type
		if path := protoPath(pass, pass.TypesInfo.Types[expr].Type); path != nil {
			pass.ReportRangef(expr, "%s passes proto by value: %v", name, path)
		}
	}

	if typ.Params != nil {
		for _, field := range typ.Params.List {
			expr := field.Type
			if path := protoPath(pass, pass.TypesInfo.Types[expr].Type); path != nil {
				pass.ReportRangef(expr, "%s passes proto by value: %v", name, path)
			}
		}
	}

	if typ.Results != nil {
		for _, field := range typ.Results.List {
			expr := field.Type
			if path := protoPath(pass, pass.TypesInfo.Types[expr].Type); path != nil {
				pass.ReportRangef(expr, "%s returns proto by value: %v", name, path)
			}
		}
	}
}

// checkCopyProtoRange checks whether a range statement
// might inadvertently copy a proto by checking whether
// any of the range variables are protos.
func checkCopyProtoRange(pass *analysis.Pass, r *ast.RangeStmt) {
	checkCopyProtoRangeVar(pass, r.Tok, r.Key)
	checkCopyProtoRangeVar(pass, r.Tok, r.Value)
}

func checkCopyProtoRangeVar(pass *analysis.Pass, rtok token.Token, e ast.Expr) {
	if e == nil {
		return
	}
	id, isID := e.(*ast.Ident)
	if isID && id.Name == "_" {
		return
	}

	var typ types.Type
	if rtok == token.DEFINE {
		if !isID {
			return
		}
		obj := pass.TypesInfo.Defs[id]
		if obj == nil {
			return
		}
		typ = obj.Type()
	} else {
		typ = pass.TypesInfo.Types[e].Type
	}

	if typ == nil {
		return
	}
	if path := protoPath(pass, typ); path != nil {
		pass.Reportf(e.Pos(), "range var %s copies proto: %v", format(pass.Fset, e), path)
	}
}

type typePath []types.Type

// String pretty-prints a typePath.
func (path typePath) String() string {
	n := len(path)
	var buf bytes.Buffer
	for i := range path {
		if i > 0 {
			_, _ = fmt.Fprint(&buf, " contains ")
		}
		// The human-readable path is in reverse order, outermost to innermost.
		_, _ = fmt.Fprint(&buf, path[n-i-1].String())
	}
	return buf.String()
}

func protoPathRhs(pass *analysis.Pass, x ast.Expr) typePath {
	if _, ok := x.(*ast.CompositeLit); ok {
		return nil
	}
	if _, ok := x.(*ast.CallExpr); ok {
		// A call may return a zero value.
		return nil
	}
	if star, ok := x.(*ast.StarExpr); ok {
		if _, ok := star.X.(*ast.CallExpr); ok {
			// A call may return a pointer to a zero value.
			return nil
		}
	}
	return protoPath(pass, pass.TypesInfo.Types[x].Type)
}

// protoPath returns a typePath describing the location of a proto value
// contained in typ. If there is no contained proto, it returns nil.
func protoPath(pass *analysis.Pass, typ types.Type) typePath {
	if typ == nil {
		return nil
	}

	for {
		atyp, ok := typ.Underlying().(*types.Array)
		if !ok {
			break
		}
		typ = atyp.Elem()
	}

	namedTyp, ok := typ.(*types.Named)
	if !ok {
		return nil
	}

	if pkg := namedTyp.Obj().Pkg(); pkg == nil || pass.ImportPackageFact(pkg, &IsGoGoPkg{}) {
		return nil
	}

	// We're only interested in the case in which the underlying
	// type is a struct.
	styp, ok := typ.Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	nfields := styp.NumFields()
	for i := 0; i < nfields; i++ {
		if styp.Field(i).Name() == "XXX_sizecache" {
			return []types.Type{typ}
		}
	}

	for i := 0; i < nfields; i++ {
		ftyp := styp.Field(i).Type()
		subpath := protoPath(pass, ftyp)
		if subpath != nil {
			return append(subpath, typ)
		}
	}

	return nil
}
