package lintexamples

import (
	Auth "a/lintexamples/internal/services/auth" // want `invalid import name Auth, use auth`
)

func dropAlias() string {
	return Auth.PackageName
}
