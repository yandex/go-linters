package lintexamples

import (
	_ "context"
	"fmt"

	apitokenauth "go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/middlewares/api_token_auth"
	"go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/services/auth"
	externalusers "go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/services/externalUsers"
)

func correctImports() {
	fmt.Println(
		apitokenauth.PackageName,
		auth.PackageName,
		externalusers.PackageName,
	)
}
