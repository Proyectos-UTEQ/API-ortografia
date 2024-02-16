package types

type PointsUserForModule struct {
	UserID    uint    `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	URLAvatar string  `json:"url_avatar"`
	Points    float64 `json:"points"`
}
