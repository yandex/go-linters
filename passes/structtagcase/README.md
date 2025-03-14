# `structtagcase`

Ensures consistent naming case in struct tags.

## What it checks

Verifies that struct tags use consistent case (snake_case or camelCase) for:
- json
- bson
- xml
- yaml
tags

## Diagnostic example

```go
// Bad - mixed cases
type User struct {
    FirstName string `json:"first_name"`
    Surname string   `json:"surName"` // want "inconsistent text case"
}

// Good - consistent
type User struct {
    FirstName string `json:"first_name"`
    Surname string   `json:"surname"`
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/structtagcase"
)

func main() {
    unitchecker.Main(structtagcase.Analyzer)
}
```

Build and run:

```
go build -o structtagcase main.go
go vet -vettool=$(which structtagcase) ./...
```
