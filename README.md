# zaplint

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

- `-key-naming-convention`: Enforce a single key naming convention (snake|kebab|camel|pascal).
- `-capitalized-message`: Enforce capitalized log messages.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.
