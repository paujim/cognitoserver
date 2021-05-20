package entities

type TokenHandler interface {
	GetTokens(username, password *string) (accessToken, refreshToken *string, err error)
	RefreshAccessToken(token *string) (accessToken, refreshToken *string, err error)
}

type UserHandler interface {
	RegisterUser(username, password *string) (sub *string, err error)
	ListUsers() (users []UserModel, err error)
}

type UserTokenHandler interface {
	TokenHandler
	UserHandler
}
