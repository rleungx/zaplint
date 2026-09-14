package zaplint

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkReplaceAny(pass *analysis.Pass, call *ast.CallExpr, pkg *types.Package) {
	t := pass.TypesInfo.TypeOf(call.Args[1])
	name := replacement(t)
	if name == "" {
		return
	}
	// Check the API in the analyzed program's zap version, not ours.
	fn, ok := pkg.Scope().Lookup(name).(*types.Func)
	if !ok {
		return
	}
	sig := fn.Type().(*types.Signature)
	if sig.Params().Len() != 2 || !types.AssignableTo(t, sig.Params().At(1).Type()) {
		return
	}
	var target ast.Node = ast.Unparen(call.Fun)
	if sel, ok := target.(*ast.SelectorExpr); ok {
		target = sel.Sel
	} else if !unshadowed(pass, call.Pos(), fn) {
		// A dot import's replacement name may be shadowed by a local binding.
		return
	}
	message := "replace zap.Any with zap." + name
	pass.Report(analysis.Diagnostic{
		Pos: call.Fun.Pos(), End: call.Fun.End(), Message: message,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message:   message,
			TextEdits: []analysis.TextEdit{{Pos: target.Pos(), End: target.End(), NewText: []byte(name)}},
		}},
	})
}

func unshadowed(pass *analysis.Pass, pos token.Pos, fn *types.Func) bool {
	for _, scope := range pass.TypesInfo.Scopes {
		if scope.Contains(pos) {
			_, object := scope.Innermost(pos).LookupParent(fn.Name(), pos)
			return object == fn
		}
	}
	return false
}

// replacement only returns constructors with the same dispatch and encoding as
// zap.Any. Assignability alone is insufficient: named types may have marshalers,
// nil errors differ from nil interfaces, and byte slices use binary encoding.
func replacement(t types.Type) string {
	if t == nil {
		return ""
	}
	t = types.Unalias(t)
	if name := scalar(t); name != "" {
		return name
	}
	switch t := t.(type) {
	case *types.Pointer:
		if name := scalar(types.Unalias(t.Elem())); name != "" {
			return name + "p"
		}
	case *types.Slice:
		elem := types.Unalias(t.Elem())
		if basic, ok := elem.(*types.Basic); ok && basic.Kind() == types.Uint8 {
			return "Binary"
		}
		if name := scalar(elem); name != "" {
			return name + "s"
		}
		if types.Identical(elem, types.Universe.Lookup("error").Type()) {
			return "Errors"
		}
	case *types.Struct:
		// Anonymous structs can promote ObjectMarshaler/ArrayMarshaler/Stringer
		// methods from embedded fields. Only method-free structs always reflect.
		if types.NewMethodSet(t).Len() == 0 {
			return "Reflect"
		}
	}
	// Arrays, named collections and interface values are deliberately retained.
	return ""
}

func scalar(t types.Type) string {
	if t, ok := t.(*types.Basic); ok {
		switch t.Kind() {
		case types.Bool:
			return "Bool"
		case types.Complex128:
			return "Complex128"
		case types.Complex64:
			return "Complex64"
		case types.Float64:
			return "Float64"
		case types.Float32:
			return "Float32"
		case types.Int:
			return "Int"
		case types.Int64:
			return "Int64"
		case types.Int32:
			return "Int32"
		case types.Int16:
			return "Int16"
		case types.Int8:
			return "Int8"
		case types.String:
			return "String"
		case types.Uint:
			return "Uint"
		case types.Uint64:
			return "Uint64"
		case types.Uint32:
			return "Uint32"
		case types.Uint16:
			return "Uint16"
		case types.Uint8:
			return "Uint8"
		case types.Uintptr:
			return "Uintptr"
		}
	}
	if named, ok := t.(*types.Named); ok && named.Obj().Pkg() != nil && named.Obj().Pkg().Path() == "time" {
		switch named.Obj().Name() {
		case "Time", "Duration":
			return named.Obj().Name()
		}
	}
	return ""
}
