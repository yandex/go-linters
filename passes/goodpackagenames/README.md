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

## Suggested fix

For an invalid import alias the analyzer suggests renaming it to the canonical name together with all its uses in the file
(`go vet -fix` applies it). The alias is dropped when the canonical name matches the imported package name.

No fix is suggested when the canonical name is not a valid identifier (e.g. a keyword) or is already taken by another
declaration visible at one of the uses: a local variable, another import, a package-level declaration or a builtin.

Package names are not fixed automatically: renaming a package affects other files and packages.

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
