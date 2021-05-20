package entities

type RegistrationRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
