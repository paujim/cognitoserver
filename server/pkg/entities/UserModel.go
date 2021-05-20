package entities

import "time"

type UserModel struct {
	Username *string    `json:"username"`
	Status   *string    `json:"status"`
	Enabled  *bool      `json:"enabled"`
	Created  *time.Time `json:"created"`
}
