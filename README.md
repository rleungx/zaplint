# zaplint

[![CI](https://img.shields.io/github/actions/workflow/status/rleungx/zaplint/go.yml?branch=main)](https://github.com/rleungx/zaplint/actions/workflows/go.yml)
[![Coverage](https://img.shields.io/codecov/c/github/rleungx/zaplint)](https://codecov.io/gh/rleungx/zaplint)
[![License](https://img.shields.io/github/license/rleungx/zaplint)](./LICENSE)

A Go static analyzer for consistent logging with [zap](https://github.com/uber-go/zap).

- Check that log messages start with an uppercase letter.
- Enforce one naming convention for structured field keys.
- Replace `zap.Any` with typed field constructors when the replacement preserves encoding behavior.

## Installation

Requires Go 1.23 or later.

```sh
go install github.com/rleungx/zaplint/cmd/zaplint@latest
```

Make sure Go's binary installation directory is on your `PATH`.

## Quick start

Run from your Go module's directory. All rules are opt-in; enable the ones you want:

```sh
zaplint \
  -capitalized-message=true \
  -key-naming-convention=snake \
  -replace-any=true \
  ./...
```

With these settings, the following call produces three diagnostics:

```go
logger.Info("request completed", zap.Any("userId", int16(42)))
```

After addressing them:

```go
logger.Info("Request completed", zap.Int16("user_id", int16(42)))
```

## Configuration

| Flag | Default | Behavior |
| --- | --- | --- |
| `-capitalized-message` | `false` | Check that messages start with an uppercase ASCII letter. |
| `-key-naming-convention` | Disabled | Check field keys using `snake`, `kebab`, `camel`, or `pascal`. |
| `-replace-any` | `false` | Suggest equivalent typed constructors for `zap.Any`. |
| `-exclude-files` | No exclusions | Skip files matching any of the comma-separated regular expressions. |
| `-fix` | `false` | Apply suggested `zap.Any` replacements. |

The rule flags `-capitalized-message` and `-replace-any` require an explicit
`true` or `false` value. Both `-replace-any=true` and `-replace-any true` work.
Use `zaplint -h` to see all available flags.

### Capitalized messages

`-capitalized-message=true` checks messages passed to `Logger` and
`SugaredLogger` logging methods, including `Log`, `Check`, and the sugared
`f`, `w`, and `ln` variants.

Messages must start with `A`–`Z`. Checks use compile-time string values, so quoted
strings, raw strings, escapes, constants, and constant concatenations are
handled consistently. Runtime strings and formatted messages with an unknown
prefix, such as `sugar.Infof("%s started", name)`, are skipped.

### Field key naming

Choose a convention with `-key-naming-convention`:

| Value | Convention | Example |
| --- | --- | --- |
| `snake` | snake_case | `user_id` |
| `kebab` | kebab-case | `user-id` |
| `camel` | camelCase | `userId` |
| `pascal` | PascalCase | `UserId` |

Checks apply to compile-time string keys in zap field constructors, including
`Object`, `Dict`, and `Namespace`, and to statically identifiable key/value
pairs in `SugaredLogger` calls such as `Infow`, `With`, and `WithLazy`.
Naming checks use ASCII letters and digits, with the separator allowed by the
selected convention.

### Typed field replacements

`-replace-any=true` suggests a constructor only when the argument's type allows
an equivalent replacement in the analyzed program's zap version. For example:

| Value type | Suggested constructor |
| --- | --- |
| `int16` | `zap.Int16` |
| `*int16` | `zap.Int16p` |
| `[]int` | `zap.Ints` |
| `[]byte` or `[]uint8` | `zap.Binary` |
| `time.Duration` | `zap.Duration` |

Byte slices use `zap.Binary` to preserve `zap.Any`'s binary encoding. Arrays,
named collections, interface values (including `error`), and custom encoding
types retain `zap.Any` when equivalence cannot be established.

## Automatic fixes

Apply safe field replacements with:

```sh
zaplint -replace-any=true -fix ./...
```

Automatic fixes apply to `zap.Any` replacements. Update message text and field
key names manually.

A run that reports diagnostics exits nonzero even when `-fix` successfully
applies edits. Rerun the check on the updated code:

```sh
zaplint -replace-any=true ./...
```

## Excluding files

For example, check field names while excluding files ending in `_generated.go`
and test files:

```sh
zaplint \
  -key-naming-convention=snake \
  -exclude-files='_generated\.go$,_test\.go$' \
  ./...
```

Patterns are Go regular expressions matched against each file's full path.
Separate patterns with commas; a match excludes the entire file from all rules.
Empty entries are ignored, so `-exclude-files=""` excludes no files.

## Using the analyzer in Go

Create an `analysis.Analyzer` with `zaplint.New` when composing a `go/analysis`
driver:

```go
analyzer := zaplint.New(&zaplint.Options{
    CapitalizedMessage:  true,
    KeyNamingConvention: zaplint.SnakeCase,
    ReplaceAny:          true,
    ExcludeFiles:        []string{`_test\.go$`},
})
```

Import `github.com/rleungx/zaplint`. Passing `nil` options leaves all rules
disabled until configured through the analyzer's flags.

## Development

```sh
make build          # Build bin/zaplint.
make test           # Prepare fixture dependencies and run the test suite.
make test-coverage  # Run tests and write coverage.out.
```

For race detection:

```sh
make test-deps
go test -race ./...
```

Tests cover diagnostics, suggested edits, and compilation and log-encoding
equivalence before and after applying replacements. Contributions are welcome
through issues and pull requests.

## License

[Apache License 2.0](./LICENSE).
