package auth_test

import (
	"testing"

	"go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/services/auth"
)

func TestConstants(t *testing.T) {
	if auth.PackageName != "auth" {
		t.Error("invalid package name")
	}
}
