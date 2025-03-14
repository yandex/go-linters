# `ctxcheck`

Ensures proper context usage in Go code.

## Analyzers

### ctxarg

Checks that `context.Context` is always the first parameter in functions.

### Diagnostic example

```go
func Good(ctx context.Context, a int) {} // OK
func Bad(a int, ctx context.Context) {} // want "context must be first argument"
```

### ctxsave

Checks that context.Context is not stored in struct fields.

### Diagnotic example

```go
type Bad struct {
    Ctx context.Context // want "context must not be saved"
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/ctxcheck"
)

func main() {
    unitchecker.Main(
        ctxcheck.CtxArgAnalyzer,
        ctxcheck.CtxSaveAnalyzer,
    )
}
```

Build and run:

```
go build -o ctxcheck main.go
go vet -vettool=$(which ctxcheck) ./...
```
