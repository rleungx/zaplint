package zaplint_test

import (
	"fmt"
	"testing"

	"github.com/rleungx/zaplint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestKeyNamingConvention(t *testing.T) {
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
