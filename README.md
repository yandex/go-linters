[![Static Badge](https://img.shields.io/badge/github-coverage-blue?logo=github)](https://yandex.github.io/go-linters/coverage)

# Go Static Analyzers Collection

This repository contains a set of static analyzers for Go code. Below is the list of available
analyzers with brief descriptions.

## Available Analyzers

1. **[copyproto](/passes/copyproto)** - Detects when protobuf messages are copied by value
2. **[deepequalproto](/passes/deepequalproto)** - Ensures protobuf messages aren't compared using
reflect.DeepEqual
3. **[goodpackagenames](/passes/goodpackagenames)** - Enforces Go naming conventions for packages and
imports
4. **[nonakedreturn](/passes/nonakedreturn)** - Prevents naked returns in functions with named
results
5. **[returnstruct](/passes/returnstruct)** - Enforces "Accept Interfaces, Return Structs" principle
6. **[structtagcase](/passes/structtagcase)** - Validates consistent casing in struct tags
7. **[remindercheck](/passes/remindercheck)** - Verifies TODO/FIXME/BUG comment formatting
8. **[ctxcheck](/passes/ctxcheck)** - Validates proper context usage (position and storage)
9. **[execinquery](/passes/execinquery)** - Detects incorrect use of Query methods for non-SELECT SQL statements

## Analyzer Middlewares

Middlewares are wrapper functions that enhance analyzers with additional capabilities. They intercept the analysis execution to add filtering or behavior modification.

### Available Middlewares

1. **Nogen** - Skips linting for generated files
   - Automatically detects generated files (using standard Go markers)
   - Suppresses all diagnostics from generated code

2. **Nolint** - Enables selective linting suppression via comments
   - Supports `//nolint:analyzername` comment directives
   - Allows developers to silence false positives
   - Adds helpful messages suggesting how to suppress reports

### Usage Examples

#### Wrapping a single analyzer

```go
import (
    "golang.org/x/tools/go/analysis/passes/nilness"
    "golang.yandex/linters/middlewares"
)

// Skip generated files for nilness analyzer
wrappedAnalyzer := middlewares.Nogen(nilness.Analyzer)

// Enable nolint comments for nilness analyzer
wrappedAnalyzer := middlewares.Nolint(nilness.Analyzer)
```

#### Combining multiple middlewares

```go
// Apply both Nogen and Nolint to an analyzer
wrappedAnalyzer := middlewares.Nolint(
    middlewares.Nogen(nilness.Analyzer),
)
```

#### Using in vettool

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"
    "golang.yandex/linters/middlewares"
    "golang.yandex/linters/passes/returnstruct"
)

func main() {
    unitchecker.Main(
        middlewares.Nolint(
            middlewares.Nogen(returnstruct.Analyzer),
        ),
    )
}
```

## Building custom vettool

Example code to create a vettool with all available analyzers:

```go
package main

import (
    "golang.org/x/tools/go/analysis/unitchecker"

    "golang.yandex/linters/passes/copyproto"
    "golang.yandex/linters/passes/ctxcheck"
    "golang.yandex/linters/passes/deepequalproto"
    "golang.yandex/linters/passes/goodpackagenames"
    "golang.yandex/linters/passes/nonakedreturn"
    "golang.yandex/linters/passes/remindercheck"
    "golang.yandex/linters/passes/returnstruct"
    "golang.yandex/linters/passes/structtagcase"
    "golang.yandex/linters/passes/execinquery"
)

func main() {
    unitchecker.Main(
        copyproto.Analyzer,
        ctxcheck.CtxArgAnalyzer,
        ctxcheck.CtxSaveAnalyzer,
        deepequalproto.Analyzer,
        goodpackagenames.Analyzer,
        nonakedreturn.Analyzer,
        remindercheck.Analyzer,
        returnstruct.Analyzer,
        structtagcase.Analyzer,
        execinquery.Analyzer,
    )
}
```

## Usage

Running via go vet:

```
go build -o yavet main.go
go vet -vettool=./yavet ./...
```

To run specific analyzer:

```
go vet -vettool=./analyzers -copyproto ./...
```
