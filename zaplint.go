package zaplint

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
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
	KeyNamingConvention string // Enforce a single key naming convention ("snake", "kebab", "camel", or "pascal").
	CapitalizedMessage  bool   // Enforce capitalized message.
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

			run(pass, opts)
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

	strVar(&opts.KeyNamingConvention, "key-naming-convention", "enforce a single key naming convention (snake|kebab|camel|pascal)")
	boolVar(&opts.CapitalizedMessage, "capitalized-message", "enforce capitalized message")
	return *fset
}

func run(pass *analysis.Pass, opts *Options) {
	visitor := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	filter := []ast.Node{(*ast.CallExpr)(nil)}

	visitor.Preorder(filter, func(node ast.Node) {
		visit(pass, opts, node)
	})
}

func visit(pass *analysis.Pass, opts *Options, node ast.Node) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	if opts.CapitalizedMessage {
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

	if opts.KeyNamingConvention != "" {
		// Check if the function being called is a zap logging function
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return
		}

		name := fn.FullName()
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
		for _, r := range key {
			if !(r == '_' || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9')) {
				return false
			}
		}
		return true
	case KebabCase:
		for _, r := range key {
			if !(r == '-' || ('a' <= r && r <= 'z') || ('0' <= r && r <= '9')) {
				return false
			}
		}
		return true
	case CamelCase:
		if len(key) == 0 || !(key[0] >= 'a' && key[0] <= 'z') {
			return false
		}
		for i := 1; i < len(key); i++ {
			if !(('a' <= key[i] && key[i] <= 'z') || ('A' <= key[i] && key[i] <= 'Z') || ('0' <= key[i] && key[i] <= '9')) {
				return false
			}
		}
		return true
	case PascalCase:
		if len(key) == 0 || !(key[0] >= 'A' && key[0] <= 'Z') {
			return false
		}
		for i := 1; i < len(key); i++ {
			if !(('a' <= key[i] && key[i] <= 'z') || ('A' <= key[i] && key[i] <= 'Z') || ('0' <= key[i] && key[i] <= '9')) {
				return false
			}
		}
		return true
	default:
		return false
	}
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
