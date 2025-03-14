# `remindercheck`

Validates format of TODO/FIXME/BUG comments.

## What it checks

Ensures reminder comments follow the pattern:
- Must be uppercase (TODO, not todo)
- Must include task ID in format TASKID-123
- Must have description

## Diagnostic example

```go
// Bad
// todo: fix this // want "must be upper case"
// TODO: implement // want "must include task id"

// Good
// TODO: TASKID-123: implement feature
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/remindercheck"
)

func main() {
    unitchecker.Main(remindercheck.Analyzer)
}
```

Build and run:

```
go build -o remindercheck main.go
go vet -vettool=$(which remindercheck) ./...
```
