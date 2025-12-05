package lintexamples

import (
	_ "context"
	"fmt"

	apitokenauth "a/lintexamples/internal/middlewares/api_token_auth"
	"a/lintexamples/internal/services/auth"
	externalusers "a/lintexamples/internal/services/externalUsers"
)

func correctImports() {
	fmt.Println(
		apitokenauth.PackageName,
		auth.PackageName,
		externalusers.PackageName,
	)
}
