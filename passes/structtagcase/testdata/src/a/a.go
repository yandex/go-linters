package a

type MiscasedUser struct {
	Name                 string `json:"name" yaml:"name"`
	Surname              string `json:"surname" yaml:"surname"`
	BankAccountNumber    string `json:"bankAccountNumber" yaml:"bank_account_number"`
	SocialSecurityNumber string `json:"social_security_number" yaml:"social_security_number"` // want `inconsistent text case in json struct tag: social_security_number`
	UsernameAlias        string `json:"username-alias" yaml:"username_alias"`                 // want `inconsistent text case in json struct tag: username-alias`
	LikedGenres          string `json:"likedGenres" yaml:"liked-genres"`                      // want `inconsistent text case in yaml struct tag: liked-genres`
	PasswordSecret       string `json:"-"`
	NothingInteresting   string
}

func InlinedStruct() {
	user := struct {
		Name                 string `json:"name"`
		Surname              string `json:"surname"`
		BankAccountNumber    string `json:"bankAccountNumber"`
		SocialSecurityNumber string `json:"social_security_number"` // want `inconsistent text case in json struct tag: social_security_number`
		UsernameAlias        string `json:"username-alias"`         // want `inconsistent text case in json struct tag: username-alias`
		LikedGenres          string `json:"likedGenres"`
	}{}

	_ = user
}

type UnknownCasing struct {
	Name                 string `bson:"name"`
	Surname              string `bson:"surname"`
	UsernameAlias        string `bson:"userName-alias"` // want `unknown casing in bson struct tag: userName-alias`
	BankAccountNumber    string `bson:"bankAccountNumber"`
	SocialSecurityNumber string `bson:"social_security_number"`
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
	UsernameAlias        string `bson:"user_name"`
	BankAccountNumber    string `bson:"account"`
	SocialSecurityNumber string `bson:"ssn"`
	LikedGenres          string `bson:"liked_genres"`
}
