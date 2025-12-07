# execinquery - a simple query string checker in Query function

> Inspired by https://github.com/1uf3/execinquery

execinquery is a linter about query string checker in Query function which reads your Go src files and
warnings it finds.

## Installation

```shell
go install github.com/lufeee/execinquery/cmd/execinquery
```

## Usage

```go
package main

import (
        "database/sql"
        "log"
)

func main() {
        db, err := sql.Open("mysql", "test:test@tcp(test:3306)/test")
        if err != nil {
                log.Fatal("Database Connect Error: ", err)
        }
        defer db.Close()

        test := "a"
        _, err = db.Query("Update * FROM hoge where id = ?", test)
        if err != nil {
                log.Fatal("Query Error: ", err)
        }

}
```

```shell
go vet -vettool=$(which execinquery) ./...

# command-line-arguments
./a.go:16:11: Use Exec instead of Query to execute `UPDATE` query
```
