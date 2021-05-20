package entities

type TokenRequest struct {
	Username     *string `form:"username"`
	Password     *string `form:"password"`
	RefreshToken *string `form:"refresh_token"`
}
