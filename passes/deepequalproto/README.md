# `deepequalproto`

Checks that protobuf messages are not compared using `reflect.DeepEqual`.

## What it checks

Detects usage of:
- `reflect.DeepEqual`
- `assert.Equal`/`assert.Equalf`
- `require.Equal`/`require.Equalf`

### Diagnostic example

```go
msg1 := &pb.Message{}
msg2 := &pb.Message{}

// Bad
if reflect.DeepEqual(msg1, msg2) { // want "avoid using reflect.DeepEqual"
    // ...
}

// Good: use proto.Equal instead
if proto.Equal(msg1, msg2) {
    // ...
}
```

## Usage

Via go vet:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/passes/deepequalproto"
)

func main() {
    unitchecker.Main(deepequalproto.Analyzer)
}
```

Build and run:

```
go build -o deepequalproto main.go
go vet -vettool=$(which deepequalproto) ./...
```
