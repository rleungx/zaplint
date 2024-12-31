package zaplint

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	SnakeCase  = "snake"
	KebabCase  = "kebab"
	CamelCase  = "camel"
	PascalCase = "pascal"
)

var errInvalidValue = errors.New("invalid value")

// Options are options for the zaplint analyzer.
type Options struct {
	CapitalizedMessage  bool     // Enforce capitalized message.
	ReplaceAny          bool     // Enforce replacing zap.Any with the appropriate type.
	KeyNamingConvention string   // Enforce a single key naming convention ("snake", "kebab", "camel", or "pascal").
	ExcludeFiles        []string // Exclude files matching the given patterns.
}

// New creates a new zaplint analyzer.
func New(opts *Options) *analysis.Analyzer {
	if opts == nil {
		opts = &Options{}
	}

	return &analysis.Analyzer{
		Name:     "zaplint",
		Doc:      "ensure consistent code style when using zap",
		Flags:    flags(opts),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			switch opts.KeyNamingConvention {
			case "", SnakeCase, KebabCase, CamelCase, PascalCase:
			default:
				return nil, fmt.Errorf("zaplint: Options.KeyNamingConvention=%s: %w", opts.KeyNamingConvention, errInvalidValue)
			}

			var regexps []*regexp.Regexp
			for _, pattern := range opts.ExcludeFiles {
				re, err := regexp.Compile(pattern)
				if err != nil {
					return nil, fmt.Errorf("zaplint: Options.ExcludeFiles=%s: %w", pattern, err)
				}
				regexps = append(regexps, re)
			}
			run(pass, opts, regexps)
			return nil, nil
		},
	}
}

func flags(opts *Options) flag.FlagSet {
	fset := flag.NewFlagSet("zaplint", flag.ContinueOnError)

	boolVar := func(value *bool, name, usage string) {
		fset.Func(name, usage, func(s string) error {
			v, err := strconv.ParseBool(s)
			*value = v
			return err
		})
	}

	strVar := func(value *string, name, usage string) {
		fset.Func(name, usage, func(s string) error {
			*value = s
			return nil
		})
	}

	strSliceVar := func(value *[]string, name, usage string) {
		fset.Func(name, usage, func(s string) error {
			*value = strings.Split(s, ",")
			return nil
		})
	}

	boolVar(&opts.CapitalizedMessage, "capitalized-message", "enforce capitalized message")
	boolVar(&opts.ReplaceAny, "replace-any", "enforce replacing zap.Any with the appropriate type")
	strVar(&opts.KeyNamingConvention, "key-naming-convention", "enforce a single key naming convention (snake|kebab|camel|pascal)")
	strSliceVar(&opts.ExcludeFiles, "exclude-files", "exclude files matching the given patterns")
	return *fset
}

func run(pass *analysis.Pass, opts *Options, regexps []*regexp.Regexp) {
	visitor := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filter := []ast.Node{(*ast.CallExpr)(nil)}

	visitor.Preorder(filter, func(node ast.Node) {
		if shouldExclude(pass.Fset.Position(node.Pos()).Filename, regexps) {
			return
		}
		visit(pass, opts, node)
	})
}

func shouldExclude(filename string, regexps []*regexp.Regexp) bool {
	for _, re := range regexps {
		if re.MatchString(filename) {
			return true
		}
	}
	return false
}

func visit(pass *analysis.Pass, opts *Options, node ast.Node) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	if opts.CapitalizedMessage {
		checkCapitalizedMessage(pass, call)
	}

	if opts.KeyNamingConvention != "" {
		checkKeyNamingConvention(pass, opts, call)
	}

	if opts.ReplaceAny {
		checkReplaceAny(pass, call)
	}
}

func checkCapitalizedMessage(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok {
		if _, ok := level[sel.Sel.Name]; ok {
			if len(call.Args) == 0 {
				return
			}

			msg, ok := call.Args[0].(*ast.BasicLit)
			if !ok || msg.Kind != token.STRING {
				return
			}

			// Trim the quotes from the message value
			msgValue := strings.Trim(msg.Value, "\"")

			if !isCapitalized(msgValue) {
				pass.Reportf(msg.Pos(), "message '%s' should be capitalized", msgValue)
			}
		}
	}
}

func checkKeyNamingConvention(pass *analysis.Pass, opts *Options, call *ast.CallExpr) {
	fn := typeutil.StaticCallee(pass.TypesInfo, call)
	if fn == nil {
		return
	}

	name := strings.TrimLeft(fn.FullName(), "vendor/")
	if _, ok := zapFields[name]; !ok {
		return
	}

	if len(call.Args) == 0 {
		return
	}

	key, ok := call.Args[0].(*ast.BasicLit)
	if !ok || key.Kind != token.STRING {
		return
	}

	// Trim the quotes from the key value
	keyValue := strings.Trim(key.Value, "\"")

	if !isValidKey(keyValue, opts.KeyNamingConvention) {
		pass.Reportf(key.Pos(), "key '%s' should be in %s", keyValue, caseMap[opts.KeyNamingConvention])
	}
}

func checkReplaceAny(pass *analysis.Pass, call *ast.CallExpr) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok && sel.Sel.Name == "Any" {
		if len(call.Args) < 2 {
			return
		}
		argType := pass.TypesInfo.TypeOf(call.Args[1])
		t := getType(argType)
		if t != "" {
			pass.Reportf(sel.Pos(), "replace zap.Any with zap.%s", t)
		}
	}
}

var caseMap = map[string]string{
	SnakeCase:  "snake_case",
	KebabCase:  "kebab-case",
	CamelCase:  "camelCase",
	PascalCase: "PascalCase",
}

func isValidKey(key, convention string) bool {
	switch convention {
	case SnakeCase:
		return isSnakeCase(key)
	case KebabCase:
		return isKebabCase(key)
	case CamelCase:
		return isCamelCase(key)
	case PascalCase:
		return isPascalCase(key)
	default:
		return false
	}
}

func isSnakeCase(key string) bool {
	for _, r := range key {
		if !(r == '_' || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9')) {
			return false
		}
	}
	return true
}

func isKebabCase(key string) bool {
	for _, r := range key {
		if !(r == '-' || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9')) {
			return false
		}
	}
	return true
}

func isCamelCase(key string) bool {
	if len(key) == 0 || !(key[0] >= 'a' && key[0] <= 'z') {
		return false
	}
	for i := 1; i < len(key); i++ {
		if !(('a' <= key[i] && key[i] <= 'z') || ('A' <= key[i] && key[i] <= 'Z') || ('0' <= key[i] && key[i] <= '9')) {
			return false
		}
	}
	return true
}

func isPascalCase(key string) bool {
	if len(key) == 0 || !(key[0] >= 'A' && key[0] <= 'Z') {
		return false
	}
	for i := 1; i < len(key); i++ {
		if !(('a' <= key[i] && key[i] <= 'z') || ('A' <= key[i] && key[i] <= 'Z') || ('0' <= key[i] && key[i] <= '9')) {
			return false
		}
	}
	return true
}

func isCapitalized(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

var zapFields = map[string]struct{}{
	"go.uber.org/zap.Binary":      {},
	"go.uber.org/zap.Bool":        {},
	"go.uber.org/zap.Boolp":       {},
	"go.uber.org/zap.ByteString":  {},
	"go.uber.org/zap.Complex128":  {},
	"go.uber.org/zap.Complex128p": {},
	"go.uber.org/zap.Complex64":   {},
	"go.uber.org/zap.Complex64p":  {},
	"go.uber.org/zap.Float64":     {},
	"go.uber.org/zap.Float64p":    {},
	"go.uber.org/zap.Float32":     {},
	"go.uber.org/zap.Float32p":    {},
	"go.uber.org/zap.Int":         {},
	"go.uber.org/zap.Intp":        {},
	"go.uber.org/zap.Int64":       {},
	"go.uber.org/zap.Int64p":      {},
	"go.uber.org/zap.Int32":       {},
	"go.uber.org/zap.Int32p":      {},
	"go.uber.org/zap.Int8":        {},
	"go.uber.org/zap.Int8p":       {},
	"go.uber.org/zap.String":      {},
	"go.uber.org/zap.Stringp":     {},
	"go.uber.org/zap.Uint":        {},
	"go.uber.org/zap.Uintp":       {},
	"go.uber.org/zap.Uint64":      {},
	"go.uber.org/zap.Uint64p":     {},
	"go.uber.org/zap.Uint32":      {},
	"go.uber.org/zap.Uint32p":     {},
	"go.uber.org/zap.Uint16":      {},
	"go.uber.org/zap.Uint16p":     {},
	"go.uber.org/zap.Uint8":       {},
	"go.uber.org/zap.Uint8p":      {},
	"go.uber.org/zap.Uintptr":     {},
	"go.uber.org/zap.Uintptrp":    {},
	"go.uber.org/zap.Reflect":     {},
	"go.uber.org/zap.Stringer":    {},
	"go.uber.org/zap.Time":        {},
	"go.uber.org/zap.Timep":       {},
	"go.uber.org/zap.Duration":    {},
	"go.uber.org/zap.Durationp":   {},
	"go.uber.org/zap.NamedError":  {},
	"go.uber.org/zap.Any":         {},

	"go.uber.org/zap.Array":        {},
	"go.uber.org/zap.Bools":        {},
	"go.uber.org/zap.ByteStrings":  {},
	"go.uber.org/zap.Complex128s":  {},
	"go.uber.org/zap.Complex64s":   {},
	"go.uber.org/zap.Durations":    {},
	"go.uber.org/zap.Float64s":     {},
	"go.uber.org/zap.Float32s":     {},
	"go.uber.org/zap.Ints":         {},
	"go.uber.org/zap.Int64s":       {},
	"go.uber.org/zap.Int32s":       {},
	"go.uber.org/zap.Int16s":       {},
	"go.uber.org/zap.Int8s":        {},
	"go.uber.org/zap.Objects":      {},
	"go.uber.org/zap.ObjectValues": {},
	"go.uber.org/zap.Strings":      {},
	"go.uber.org/zap.Stringers":    {},
	"go.uber.org/zap.Times":        {},
	"go.uber.org/zap.Uints":        {},
	"go.uber.org/zap.Uint64s":      {},
	"go.uber.org/zap.Uint32s":      {},
	"go.uber.org/zap.Uint16s":      {},
	"go.uber.org/zap.Uint8s":       {},
	"go.uber.org/zap.Uintptrs":     {},
	"go.uber.org/zap.Errors":       {},
}

var level = map[string]struct{}{
	"Debug":  {},
	"Info":   {},
	"Warn":   {},
	"Error":  {},
	"DPanic": {},
	"Panic":  {},
	"Fatal":  {},
}

func getType(t types.Type) string {
	switch t.String() {
	case "bool":
		return "Bool"
	case "*bool":
		return "Boolp"
	case "complex128":
		return "Complex128"
	case "*complex128":
		return "Complex128p"
	case "complex64":
		return "Complex64"
	case "*complex64":
		return "Complex64p"
	case "float64":
		return "Float64"
	case "*float64":
		return "Float64p"
	case "float32":
		return "Float32"
	case "*float32":
		return "Float32p"
	case "int":
		return "Int"
	case "*int":
		return "Intp"
	case "int64":
		return "Int64"
	case "*int64":
		return "Int64p"
	case "int32":
		return "Int32"
	case "*int32":
		return "Int32p"
	case "int8":
		return "Int8"
	case "*int8":
		return "Int8p"
	case "string":
		return "String"
	case "*string":
		return "Stringp"
	case "uint":
		return "Uint"
	case "*uint":
		return "Uintp"
	case "uint64":
		return "Uint64"
	case "*uint64":
		return "Uint64p"
	case "uint32":
		return "Uint32"
	case "*uint32":
		return "Uint32p"
	case "uint16":
		return "Uint16"
	case "*uint16":
		return "Uint16p"
	case "uint8":
		return "Uint8"
	case "*uint8":
		return "Uint8p"
	case "uintptr":
		return "Uintptr"
	case "*uintptr":
		return "Uintptrp"
	case "time.Duration":
		return "Duration"
	case "*time.Duration":
		return "Durationp"
	case "time.Time":
		return "Time"
	case "*time.Time":
		return "Timep"
	case "error":
		return "NamedError"
	default:
		if strings.HasPrefix(t.String(), "struct") {
			return "Reflect"
		}
		switch t.Underlying().(type) {
		case *types.Array, *types.Slice:
			elemType := t.Underlying().(interface {
				Elem() types.Type
			}).Elem().String()
			switch elemType {
			case "bool":
				return "Bools"
			case "[]byte":
				return "ByteStrings"
			case "complex128":
				return "Complex128s"
			case "complex64":
				return "Complex64s"
			case "time.Duration":
				return "Durations"
			case "float64":
				return "Float64s"
			case "float32":
				return "Float32s"
			case "int":
				return "Ints"
			case "int64":
				return "Int64s"
			case "int32":
				return "Int32s"
			case "int16":
				return "Int16s"
			case "int8":
				return "Int8s"
			case "string":
				return "Strings"
			case "time.Time":
				return "Times"
			case "uint":
				return "Uints"
			case "uint64":
				return "Uint64s"
			case "uint32":
				return "Uint32s"
			case "uint16":
				return "Uint16s"
			case "uint8":
				return "Uint8s"
			case "uintptr":
				return "Uintptrs"
			case "error":
				return "Errors"
			}
		}
	}
	return ""
}
