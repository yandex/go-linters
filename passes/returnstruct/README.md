# `returnstruct`

Enforces the "Accept Interfaces, Return Structs" principle for return values.

## What it checks

Detects when:
- Function returns an interface instead of concrete type
- Exception: built-in `error` interface is allowed

## Diagnostic example

```go
// Bad
func Bad() io.Closer { // want "function must return concrete type"
    return &User{}
}

// Good
func Good() *User {
    return &User{}
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/returnstruct"
)

func main() {
    unitchecker.Main(returnstruct.Analyzer)
}
```

Build and run:

```
go build -o returnstruct main.go
go vet -vettool=$(which returnstruct) ./...
```

