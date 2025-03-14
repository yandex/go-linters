package externalUsers_test // want `invalid package name externalUsers_test, use externalusers_test`

import (
	"testing"

	externalusers "go-linters/passes/goodpackagenames/testdata/src/a/lintexamples/internal/services/externalUsers"
)

func TestConstants(t *testing.T) {
	if externalusers.PackageName != "externalUsers" {
		t.Error("invalid package name")
	}
}
