package zaplint

import (
	"go/ast"
	"go/constant"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

const zapPackage = "go.uber.org/zap"

// GOPATH type information retains vendor prefixes; module mode does not.
// Strip only complete vendor path components, never a set of characters.
func packagePath(pkg *types.Package) string {
	if pkg == nil {
		return ""
	}
	path := pkg.Path()
	if i := strings.LastIndex(path, "/vendor/"); i >= 0 {
		return path[i+len("/vendor/"):]
	}
	return strings.TrimPrefix(path, "vendor/")
}

func checkCall(pass *analysis.Pass, opts *Options, call *ast.CallExpr) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil || packagePath(fn.Pkg()) != zapPackage {
		return
	}
	sig := fn.Type().(*types.Signature)
	args := call.Args
	if sig.Recv() == nil {
		if opts.KeyNamingConvention != "" && isFieldConstructor(sig) && len(args) > 0 {
			checkKey(pass, opts.KeyNamingConvention, args[0])
		}
		if opts.ReplaceAny && fn.Name() == "Any" && len(args) == 2 {
			checkReplaceAny(pass, call, fn.Pkg())
		}
		return
	}

	receiver := namedType(sig.Recv().Type())
	if receiver == nil || (receiver.Obj().Name() != "Logger" && receiver.Obj().Name() != "SugaredLogger") {
		return
	}
	// In a method expression the explicit receiver precedes the parameters.
	if sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr); ok {
		if selection := pass.TypesInfo.Selections[sel]; selection != nil && selection.Kind() == types.MethodExpr {
			if len(args) == 0 {
				return
			}
			args = args[1:]
		}
	}
	sugared := receiver.Obj().Name() == "SugaredLogger"
	message, keys := loggerArguments(fn.Name(), sugared)
	if opts.CapitalizedMessage && message >= 0 && message < len(args) {
		checkMessage(pass, fn.Name(), sugared, args[message:])
	}
	if opts.KeyNamingConvention != "" && keys >= 0 && keys < len(args) {
		checkSugaredKeys(pass, opts.KeyNamingConvention, args[keys:])
	}
}

func namedType(t types.Type) *types.Named {
	t = types.Unalias(t)
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}
	named, _ := t.(*types.Named)
	return named
}

func isField(t types.Type) bool {
	named, ok := types.Unalias(t).(*types.Named)
	return ok && packagePath(named.Obj().Pkg()) == zapPackage+"/zapcore" && named.Obj().Name() == "Field"
}

// All keyed zap field constructors take a string first and return a Field.
// Use their actual signatures so generic and newly added constructors work too.
func isFieldConstructor(sig *types.Signature) bool {
	return sig.Params().Len() > 0 && types.Identical(sig.Params().At(0).Type(), types.Typ[types.String]) &&
		sig.Results().Len() == 1 && isField(sig.Results().At(0).Type())
}

// Indices are relative to parameters, excluding an explicit method receiver.
func loggerArguments(name string, sugared bool) (message, keys int) {
	if name == "Check" || name == "Log" {
		return 1, -1
	}
	if sugared {
		switch name {
		case "With", "WithLazy":
			return -1, 0
		case "Logf", "Logln":
			return 1, -1
		case "Logw":
			return 1, 2
		}
	}
	base := name
	if sugared {
		for _, suffix := range []string{"ln", "f", "w"} {
			if strings.HasSuffix(base, suffix) {
				base = strings.TrimSuffix(base, suffix)
				break
			}
		}
	}
	switch base {
	case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
		if sugared && strings.HasSuffix(name, "w") {
			return 0, 1
		}
		return 0, -1
	}
	return -1, -1
}

func constantString(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	value := pass.TypesInfo.Types[expr].Value
	if value == nil || value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(value), true
}

func checkMessage(pass *analysis.Pass, method string, sugared bool, args []ast.Expr) {
	msg, ok := constantString(pass, args[0])
	if !ok {
		return
	}
	if sugared {
		if !types.Identical(pass.TypesInfo.TypeOf(args[0]), types.Typ[types.String]) {
			// Sprint can invoke String/Format methods on a defined string type.
			return
		}
		// A formatted prefix or empty Sprint prefix can depend on later args.
		// Without a known prefix there is no sound capitalization diagnostic.
		if strings.HasSuffix(method, "f") && len(args) > 1 && strings.HasPrefix(msg, "%") {
			return
		}
		if msg == "" && len(args) > 1 && !strings.HasSuffix(method, "w") && method != "Check" {
			return
		}
	}
	if !isCapitalized(msg) {
		pass.Reportf(args[0].Pos(), "message '%s' should be capitalized", msg)
	}
}

func checkKey(pass *analysis.Pass, convention string, expr ast.Expr) {
	if key, ok := constantString(pass, expr); ok && !isValidKey(key, convention) {
		pass.Reportf(expr.Pos(), "key '%s' should be in %s", key, caseMap[convention])
	}
}

func checkSugaredKeys(pass *analysis.Pass, convention string, args []ast.Expr) {
	errorType := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	for len(args) > 0 {
		t := pass.TypesInfo.TypeOf(args[0])
		if t == nil {
			return
		}
		// Field and concrete error arguments consume one slot in sweetenFields.
		if isField(t) || (!types.IsInterface(t) && types.Implements(t, errorType)) {
			args = args[1:]
			continue
		}
		// An interface may hold a Field, error or key. Its runtime value controls
		// the alignment of every following key/value pair, including nil errors.
		if types.IsInterface(t) || len(args) < 2 {
			return
		}
		if types.Identical(t, types.Typ[types.String]) {
			checkKey(pass, convention, args[0])
		}
		args = args[2:]
	}
}
