package zaplint_test

import (
	"testing"

	"github.com/rleungx/zaplint"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEmptyExclusions(t *testing.T) {
	t.Parallel()
	for _, flag := range []string{"", ",", "does_not_match,", ",does_not_match", " , "} {
		t.Run(flag, func(t *testing.T) {
			a := zaplint.New(&zaplint.Options{KeyNamingConvention: zaplint.SnakeCase})
			if err := a.Flags.Set("exclude-files", flag); err != nil {
				t.Fatal(err)
			}
			analysistest.Run(t, analysistest.TestData(), a, "key_naming_convention/snake")
		})
	}
	a := zaplint.New(&zaplint.Options{KeyNamingConvention: zaplint.SnakeCase, ExcludeFiles: []string{""}})
	analysistest.Run(t, analysistest.TestData(), a, "key_naming_convention/snake")
}

func TestInvalidOptions(t *testing.T) {
	t.Parallel()
	for _, opts := range []zaplint.Options{
		{KeyNamingConvention: "invalid"},
		{ExcludeFiles: []string{"["}},
	} {
		a := zaplint.New(&opts)
		if _, err := a.Run(&analysis.Pass{}); err == nil {
			t.Fatalf("expected invalid options to fail: %+v", opts)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	t.Parallel()
	for _, opts := range []*zaplint.Options{nil, {}} {
		analysistest.Run(t, analysistest.TestData(), zaplint.New(opts), "disabled")
	}
}

func TestCallAliasesAndFlags(t *testing.T) {
	t.Parallel()
	a := zaplint.New(nil)
	if err := a.Flags.Parse([]string{"-capitalized-message", "true", "-replace-any", "true", "-key-naming-convention", "snake"}); err != nil {
		t.Fatal(err)
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), a, "call_aliases")
}
