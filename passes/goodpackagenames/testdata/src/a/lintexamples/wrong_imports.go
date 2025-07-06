package lintexamples

import (
	"fmt"

	api_token_auth "a/lintexamples/internal/middlewares/api_token_auth" // want `invalid import name api_token_auth, use apitokenauth`
	externalUsers "a/lintexamples/internal/services/externalUsers"      // want `invalid import name externalUsers, use externalusers`
)

func wrongImports() {
	fmt.Println(
		api_token_auth.PackageName,
		externalUsers.PackageName,
	)
}
