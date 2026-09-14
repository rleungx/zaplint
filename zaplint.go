package zaplint

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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
				if strings.TrimSpace(pattern) == "" {
					continue
				}
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
			if err != nil {
				return err
			}
			*value = v
			return nil
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
	if !opts.CapitalizedMessage && !opts.ReplaceAny && opts.KeyNamingConvention == "" {
		return
	}
	visitor := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	// Exclusion is a file property: prune the entire file before visiting calls.
	filter := []ast.Node{(*ast.File)(nil), (*ast.CallExpr)(nil)}
	visitor.Nodes(filter, func(node ast.Node, push bool) bool {
		if !push {
			return false
		}
		switch node := node.(type) {
		case *ast.File:
			return !shouldExclude(pass.Fset.File(node.Pos()).Name(), regexps)
		case *ast.CallExpr:
			checkCall(pass, opts, node)
		}
		return true
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
