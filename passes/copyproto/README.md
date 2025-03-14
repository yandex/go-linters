# `copyproto`

Checks that protocol buffer messages are not copied by value.

## What it checks

This linter detects situations where protobuf messages are being copied, which can lead to
performance issues and unintended behavior. It checks:

- Variable assignments
- Function calls with proto arguments
- Return statements
- Composite literals
- Range loops
- Variable declarations

## Diagnostic example

```go
type Msg struct { XXX_sizecache int32 }

// Bad - returns proto by value
func X() Msg {
    return Msg{} // want "returns proto by value"
}

// Bad - passes proto by value
func Y(Msg) {}

func Z() {
    var m Msg
    _ = m // want "assignment copies proto value"
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/copyproto"
)

func main() {
    unitchecker.Main(copyproto.Analyzer)
}
```

Build and run:

```
go build -o copyproto main.go
go vet -vettool=$(which copyproto) ./...
```

### Recommendations

Always pass protobuf messages by pointer (*Msg) rather than by value.
