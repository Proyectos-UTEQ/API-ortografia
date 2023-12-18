package types

type UserAPI struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
