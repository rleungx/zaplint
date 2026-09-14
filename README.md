# zaplint

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rleungx/zaplint/go.yml)
![Codecov](https://img.shields.io/codecov/c/github/rleungx/zaplint)
![GitHub License](https://img.shields.io/github/license/rleungx/zaplint)

`zaplint` is a Go static analysis tool that ensures consistent code style when using the `zap` logging library.

## Features

- Enforce capitalized log messages.
- Enforce replacing `zap.Any` with the appropriate type.
- Enforce a single key naming convention: snake_case, kebab-case, camelCase, or PascalCase.
- Exclude specified files or patterns from analysis.

## Installation

To install `zaplint`, use the following command:

```sh
go install github.com/rleungx/zaplint/cmd/zaplint@latest
```

## Usage

You can run `zaplint` through the following command: 

```sh
zaplint -key-naming-convention kebab -capitalized-message true -replace-any true ./...
```

## Configuration

You can configure `zaplint` using the following flags:

- `-capitalized-message`: Enforce capitalized log messages.
- `-replace-any`: Enforce replacing `zap.Any` with the appropriate type.
- `-key-naming-convention`: Enforce a single key naming convention (`snake`|`kebab`|`camel`|`pascal`).
- `-exclude-files`: Exclude files matching the given patterns (comma-separated).

Messages and keys are checked using their compile-time string values, including
raw strings, escapes, constants, and constant concatenations. Message checks
cover `Logger` and `SugaredLogger` logging methods, including `Log`, `Check`,
and the `f`, `w`, and `ln` variants. A formatted message with a dynamic prefix
is skipped. Key checks cover keyed zap field constructors and statically
identifiable key/value pairs in `SugaredLogger` calls.

`-replace-any` only suggests replacements that preserve zap's encoding behavior.
Arrays, named collections, interface values (including `error`), and values with
custom encoding are retained when equivalence cannot be established. Byte slices
use `zap.Binary`. Safe suggestions can be applied with the standard `-fix` flag.
A run that reports diagnostics still exits nonzero after applying fixes; rerun
the command to check the updated code.

File exclusion patterns are regular expressions. Empty entries are ignored;
`-exclude-files=""` excludes no files. All rules are opt-in and retain their
existing ASCII naming and capitalization conventions.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.
