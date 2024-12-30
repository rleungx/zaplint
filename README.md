# zaplint

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/rleungx/zaplint/go.yml)
![Codecov](https://img.shields.io/codecov/c/github/rleungx/zaplint)
![GitHub License](https://img.shields.io/github/license/rleungx/zaplint)

`zaplint` is a Go static analysis tool that ensures consistent code style when using the `zap` logging library.

## Features

- Enforce a single key naming convention: snake_case, kebab-case, camelCase, or PascalCase.
- Enforce capitalized log messages.

## Installation

To install `zaplint`, use the following command:

```sh
go get github.com/rleungx/zaplint
```

## Usage

You can run `zaplint` through the following command: 

```sh
zaplint -key-naming-convention kebab -capitalized-message true ./...
```

## Configuration

You can configure `zaplint` using the following flags:

- `-capitalized-message`: Enforce capitalized log messages.
- `-key-naming-convention`: Enforce a single key naming convention (`snake`|`kebab`|`camel`|`pascal`).
- `-exclude-files`: Exclude files matching the given patterns (comma-separated).

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.
