package force_camel

type MiscasedUser struct {
	Name                 string `json:"name" yaml:"name"`
	Surname              string `json:"surname" yaml:"surname"`
	BankAccountNumber    string `json:"bankAccountNumber" yaml:"bank_account_number"` // want `yaml struct tag must be in camel case: bank_account_number`
	SocialSecurityNumber string `json:"social_security_number" yaml:"social_security_number"` // want `json struct tag must be in camel case: social_security_number` `yaml struct tag must be in camel case: social_security_number`
	UsernameAlias        string `json:"username-alias" yaml:"username_alias"` // want `json struct tag must be in camel case: username-alias` `yaml struct tag must be in camel case: username_alias`
	LikedGenres          string `json:"likedGenres" yaml:"liked-genres"` // want `yaml struct tag must be in camel case: liked-genres`
	PasswordSecret       string `json:"-"`
	NothingInteresting   string
}

func InlinedStruct() {
	user := struct {
		Name                 string `json:"name"`
		Surname              string `json:"surname"`
		BankAccountNumber    string `json:"bankAccountNumber"`
		SocialSecurityNumber string `json:"social_security_number"` // want `json struct tag must be in camel case: social_security_number`
		UsernameAlias        string `json:"username-alias"` // want `json struct tag must be in camel case: username-alias`
		LikedGenres          string `json:"likedGenres"`
	}{}

	_ = user
}

type UnknownCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName-alias"`    // want `bson struct tag must be in camel case: userName-alias`
	BankAccountNumber    string `bson:"bankAccountNumber"`
	SocialSecurityNumber string `bson:"social_security_number"` // want `bson struct tag must be in camel case: social_security_number`
	LikedGenres          string `bson:"likedGenres"`
	PasswordSecret       string `bson:"-"`
}

type AcceptableCamelCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName"`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"likedGenres"`
}

type AcceptableSnakeCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"user_name"` // want `bson struct tag must be in camel case: user_name`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"liked_genres"` // want `bson struct tag must be in camel case: liked_genres`
}
