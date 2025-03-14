# `nonakedreturn`

Checks that functions with named results don't use naked returns.

## What it checks

Detects cases when:
- Function has named return values
- Return statement is used without explicit values

## Diagnostic example

```go
// Bad
func bad() (x int) {
    x = 5
    return // want "Naked return"
}

// Good
func good() (x int) {
    return 5
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/nonakedreturn"
)

func main() {
    unitchecker.Main(nonakedreturn.Analyzer)
}
```

Build and run:

```
go build -o nonakedreturn main.go
go vet -vettool=$(which nonakedreturn) ./...
```
