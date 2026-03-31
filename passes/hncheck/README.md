# `hncheck`

Checks for variables/constants/types names for Hungarian notation usage.

## Diagnostic example

```go
// Bad
var iCount int
var sName string
var fPrice float32

// Good
var count int
var name string
var price float32
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/hncheck"
)

func main() {
    unitchecker.Main(hncheck.Analyzer)
}
```

Build and run:

```
go build -o hncheck main.go
go vet -vettool=$(which hncheck) ./...
```
