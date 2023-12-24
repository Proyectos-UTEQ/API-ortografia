package types

type UserAPI struct {
	ID           uint   `json:"id"`
	FirstName    string `json:"first_name" validate:"required,min=6,max=100"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password,omitempty" validate:"required,min=6,max=100"`
	BirthDate    string `json:"birth_date"`
	PointsEarned int    `json:"points_earned"`
	Whatsapp     string `json:"whatsapp"`
	Telegram     string `json:"telegram"`
	URLAvatar    string `json:"url_avatar"`
	Status       string `json:"status"`
	TypeUser     string `json:"type_user" validate:"required,oneof=student teacher admin"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
