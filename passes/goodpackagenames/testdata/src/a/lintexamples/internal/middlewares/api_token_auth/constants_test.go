package api_token_auth_test // want `invalid package name api_token_auth_test, use apitokenauth_test`

import (
	"testing"

	apitokenauth "go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/middlewares/api_token_auth"
)

func TestConstants(t *testing.T) {
	if apitokenauth.PackageName != "api_token_auth" {
		t.Error("invalid package name")
	}
}
