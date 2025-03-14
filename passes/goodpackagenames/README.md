# `goodpackagenames`

Checks that package and import names follow Go naming conventions.

## What it checks

Verifies that:
- Package names are lowercase with no underscores or mixedCaps
- Import aliases follow the same naming rules

## Diagnostic example

```go
// Bad
package token_auth  // want "invalid package name"

// Good
package tokenauth

// Bad import
import tokenAuth "path"  // want "invalid import name"

// Good import
import tokenauth "path"
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/goodpackagenames"
)

func main() {
    unitchecker.Main(goodpackagenames.Analyzer)
}
```

Build and run:

```
go build -o goodpackagenames main.go
go vet -vettool=$(which goodpackagenames) ./...
```
