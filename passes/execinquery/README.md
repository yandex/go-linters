# execinquery

`execinquery` is a linter about query string checker in Query function which reads your Go src files and
warnings it finds.

## Features

- Detects when `Query`, `QueryRow`, `QueryContext`, or `QueryRowContext` are used with non-SELECT queries
- Suggests using `Exec` or `ExecContext` instead for INSERT, UPDATE, DELETE queries
- Supports PostgreSQL `RETURNING` clauses (queries with RETURNING are allowed to use Query/QueryRow)
- Allows transaction control statements (`BEGIN`) to use Query methods (for compatibility)
- Handles SQL comments (single-line `--` and multi-line `/* */`)

> # Disclaimer
>
> This is a fork of the original linter repository [execinquery](https://github.com/1uf3/execinquery).
>
> - Retains core functionality with possible modifications.
> - All original code made before first commit in this repository distributes under MIT license.
> - Changes here may not sync with the original.

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/execinquery"
)

func main() {
    unitchecker.Main(execinquery.Analyzer)
}
```

Build and run:

```
go build -o execinquery main.go
go vet -vettool=$(which execinquery) ./...
```
