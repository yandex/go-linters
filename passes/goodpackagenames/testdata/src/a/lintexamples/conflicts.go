package lintexamples

import (
	API_Token_Auth "a/lintexamples/internal/middlewares/api_token_auth" // want `invalid import name API_Token_Auth, use apitokenauth`
	Len "a/lintexamples/internal/services/auth"                         // want `invalid import name Len, use len`
	Type "a/lintexamples/internal/services/auth"                        // want `invalid import name Type, use type`
	Ext_Users "a/lintexamples/internal/services/externalUsers"          // want `invalid import name Ext_Users, use extusers`
)

// no fix: the new name is taken by a package-level declaration
var extusers = "clash"

// no fix: the new name is taken by a local variable at one of the uses
func localClash() string {
	apitokenauth := "clash"

	return API_Token_Auth.PackageName + apitokenauth
}

func packageClash() string {
	return Ext_Users.PackageName + extusers
}

// no fix: the new name would shadow a builtin
func builtinClash() int {
	return len(Len.PackageName)
}

// no fix: the new name is a keyword
func keywordClash() string {
	return Type.PackageName
}
