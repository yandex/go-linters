package force_snake

type MiscasedUser struct {
	Name                 string `json:"name" yaml:"name"`
	Surname              string `json:"surname" yaml:"surname"`
	BankAccountNumber    string `json:"bankAccountNumber" yaml:"bank_account_number"` // want `json struct tag must be in snake case: bankAccountNumber`
	SocialSecurityNumber string `json:"social_security_number" yaml:"social_security_number"`
	UsernameAlias        string `json:"username-alias" yaml:"username_alias"` // want `json struct tag must be in snake case: username-alias`
	LikedGenres          string `json:"likedGenres" yaml:"liked-genres"`      // want `json struct tag must be in snake case: likedGenres` `yaml struct tag must be in snake case: liked-genres`
	PasswordSecret       string `json:"-"`
	NothingInteresting   string
}

func InlinedStruct() {
	user := struct {
		Name                 string `json:"name"`
		Surname              string `json:"surname"`
		BankAccountNumber    string `json:"bankAccountNumber"` // want `json struct tag must be in snake case: bankAccountNumber`
		SocialSecurityNumber string `json:"social_security_number"`
		UsernameAlias        string `json:"username-alias"` // want `json struct tag must be in snake case: username-alias`
		LikedGenres          string `json:"likedGenres"`    // want `json struct tag must be in snake case: likedGenres`
	}{}

	_ = user
}

type UnknownCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName-alias"`    // want `bson struct tag must be in snake case: userName-alias`
	BankAccountNumber    string `bson:"bankAccountNumber"` // want `bson struct tag must be in snake case: bankAccountNumber`
	SocialSecurityNumber string `bson:"social_security_number"`
	LikedGenres          string `bson:"likedGenres"` // want `bson struct tag must be in snake case: likedGenres`
	PasswordSecret       string `bson:"-"`
}

type AcceptableCamelCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName"` // want `bson struct tag must be in snake case: userName`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"likedGenres"` // want `bson struct tag must be in snake case: likedGenres`
}

type AcceptableSnakeCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"user_name"`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"liked_genres"`
}
