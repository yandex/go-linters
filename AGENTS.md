# Go-Linters Agent Guide

This guide helps agents work effectively in the Yandex Go linters repository.

## Project Overview

This is a collection of static analyzers for Go code, built using `golang.org/x/tools/go/analysis`. Each analyzer is a separate pass that can be used individually or combined into a vet tool.

## Essential Commands

### Building
```bash
go build ./...                    # Build all packages
go mod tidy                       # Clean dependencies
```

### Testing
```bash
go test -race ./...               # Run all tests with race detection
go test -race ./passes/...        # Run tests for specific analyzer passes
go test -run TestSpecific ./passes/analyzername  # Run specific test
```

### Coverage
```bash
go test -race ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
```

## Code Organization

### Directory Structure
- `passes/[analyzer]/` - Individual analyzer packages
  - `analyzer.go` - Main analyzer implementation
  - `analyzer_test.go` - Tests using `analysistest`
  - `testdata/src/[case]/` - Test input files with `// want` comments
- `internal/` - Shared utilities
  - `lintutils/` - Common linting utilities
  - `nolint/` // Handles `//nolint:` directives
  - `nogen/` - Filters generated files
- `.github/workflows/` - CI configuration

### Analyzer Pattern
Each analyzer follows this structure:
```go
var Analyzer = &analysis.Analyzer{
    Name: Name,
    Doc:  "Description of what it does",
    Run:  run,
    Flags: flags,  // Optional
    Requires: []*analysis.Analyzer{
        nolint.Analyzer,
        nogen.Analyzer,
    },
}
```

## Testing Approach

### Testdata Structure
- Test cases go in `testdata/src/[testcase]/`
- Use `// want` comments inline to mark expected diagnostics
- Multiple test cases for different scenarios (error cases, success cases)

### Standard Test Pattern
```go
func Test(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, Analyzer, "testcase_package")
}
```

### Testing with Flags
Set package-level flags and reset after use:
```go
flagForceCasing = casingSnake
defer func() { flagForceCasing = casingUnknown }()
```

## Code Patterns and Conventions

### Shared Dependencies
All analyzers use these shared utilities:
- `golang.org/x/tools/go/analysis` - Core analysis framework
- `golang.org/x/tools/go/ast/inspector` - AST traversal
- `internal/nolint` - Handle `//nolint:` comments
- `internal/nogen` - Filter generated files
- `internal/lintutils` - Common helpers

### Analyzer Implementation Pattern
```go
func run(pass *analysis.Pass) (any, error) {
    nogenFiles := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files)
    nolintIndex := lintutils.ResultOf(pass, nolint.Name).(*nolint.Index)
    
    ins := inspector.New(nogenFiles.List())
    
    nodeFilter := []ast.Node{
        (*ast.TargetNodeType)(nil),
    }
    
    return nil, nil
}
```

### Flag Handling
Use package-level flag variables with proper initialization:
```go
var (
    flags flag.FlagSet
    flagForceCasing stringCasing
)

func init() {
    flags.Var(&flagForceCasing, "force-casing", "description")
}
```

## Important Gotchas

### Generated File Filtering
Always use `nogen.Analyzer` to filter generated files:
```go
nogenFiles := lintutils.ResultOf(pass, nogen.Name).(*nogen.Files)
ins := inspector.New(nogenFiles.List())
```

### Nolint Handling
Always respect `//nolint:` directives:
```go
nolintIndex := lintutils.ResultOf(pass, nolint.Name).(*nolint.Index)
nolintNodes := nolintIndex.ForLinter(Name)
if nolintNodes.Excluded(node) {
    continue
}
```

### Test File Naming
- Test files ending in `_test.go` should be skipped in analysis

### Inspector Usage
Use push=false filtering to avoid processing nodes twice:
```go
ins.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
    if !push {
        return false
    }
    // Process node
    return true
})
```

## Adding New Analyzers

1. Create directory `passes/analyzername/`
2. Implement analyzer following the pattern above
3. Add testdata in `testdata/src/`
4. Add `analysistest.Run()` tests
5. Update README.md with analyzer description
6. Add to example vettool in main.go

## CI/CD

- Tests run on two latest minor Go releases across Ubuntu, macOS, Windows
- Coverage is auto-generated and deployed to GitHub Pages
- Use `-race` flag in all CI tests
