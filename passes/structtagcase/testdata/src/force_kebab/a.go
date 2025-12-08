package force_kebab

type MiscasedUser struct {
	Name                 string `json:"name" yaml:"name"`
	Surname              string `json:"surname" yaml:"surname"`
	BankAccountNumber    string `json:"bankAccountNumber" yaml:"bank_account_number"` // want `json struct tag must be in kebab case: bankAccountNumber` `yaml struct tag must be in kebab case: bank_account_number`
	SocialSecurityNumber string `json:"social_security_number" yaml:"social_security_number"` // want `json struct tag must be in kebab case: social_security_number` `yaml struct tag must be in kebab case: social_security_number`
	UsernameAlias        string `json:"username-alias" yaml:"username_alias"` // want `yaml struct tag must be in kebab case: username_alias`
	LikedGenres          string `json:"likedGenres" yaml:"liked-genres"` // want `json struct tag must be in kebab case: likedGenres`
	PasswordSecret       string `json:"-"`
	NothingInteresting   string
}

func InlinedStruct() {
	user := struct {
		Name                 string `json:"name"`
		Surname              string `json:"surname"`
		BankAccountNumber    string `json:"bankAccountNumber"` // want `json struct tag must be in kebab case: bankAccountNumber`
		SocialSecurityNumber string `json:"social_security_number"` // want `json struct tag must be in kebab case: social_security_number`
		UsernameAlias        string `json:"username-alias"`
		LikedGenres          string `json:"likedGenres"` // want `json struct tag must be in kebab case: likedGenres`
	}{}

	_ = user
}

type UnknownCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName-alias"`    // want `bson struct tag must be in kebab case: userName-alias`
	BankAccountNumber    string `bson:"bankAccountNumber"` // want `bson struct tag must be in kebab case: bankAccountNumber`
	SocialSecurityNumber string `bson:"social_security_number"` // want `bson struct tag must be in kebab case: social_security_number`
	LikedGenres          string `bson:"likedGenres"` // want `bson struct tag must be in kebab case: likedGenres`
	PasswordSecret       string `bson:"-"`
}

type AcceptableCamelCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName"` // want `bson struct tag must be in kebab case: userName`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"likedGenres"` // want `bson struct tag must be in kebab case: likedGenres`
}

type AcceptableSnakeCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"user_name"` // want `bson struct tag must be in kebab case: user_name`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"liked_genres"` // want `bson struct tag must be in kebab case: liked_genres`
}
