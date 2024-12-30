package zaplint_test

import (
	"fmt"
	"testing"

	"github.com/rleungx/zaplint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestKeyNamingConvention(t *testing.T) {
	t.Parallel()
	conventions := map[string]string{
		"snake":  zaplint.SnakeCase,
		"camel":  zaplint.CamelCase,
		"kebab":  zaplint.KebabCase,
		"pascal": zaplint.PascalCase,
	}

	for name, convention := range conventions {
		t.Run(name, func(t *testing.T) {
			opts := &zaplint.Options{KeyNamingConvention: convention}
			analyzer := zaplint.New(opts)
			analysistest.Run(t, analysistest.TestData(), analyzer, fmt.Sprintf("key_naming_convention/%s", name))
		})
	}
}

func TestCapitalizedMessage(t *testing.T) {
	t.Parallel()
	opts := &zaplint.Options{CapitalizedMessage: true}
	analyzer := zaplint.New(opts)
	analysistest.Run(t, analysistest.TestData(), analyzer, "capitalized_message")
}

func TestExcludeFiles(t *testing.T) {
	t.Parallel()
	opts := &zaplint.Options{
		KeyNamingConvention: zaplint.SnakeCase,
		ExcludeFiles:        []string{`excluded_file_foo`, `excluded_file_bar`},
	}
	analyzer := zaplint.New(opts)
	analysistest.Run(t, analysistest.TestData(), analyzer, "exclude_files")

	opts = &zaplint.Options{
		KeyNamingConvention: zaplint.SnakeCase,
		ExcludeFiles:        []string{`excluded_file*`},
	}
	analyzer = zaplint.New(opts)
	analysistest.Run(t, analysistest.TestData(), analyzer, "exclude_files")
}
