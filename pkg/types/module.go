package types

// Representacion de un modulo para el frontend
type Module struct {
	ID               uint    `json:"id"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	CreateBy         UserAPI `json:"create_by" validate:"-"`
	Code             string  `json:"code"`
	Title            string  `json:"title" validate:"required,min=3,max=100"`
	ShortDescription string  `json:"short_description" validate:"required,min=3,max=100"`
	TextRoot         string  `json:"text_root" `
	ImgBackURL       string  `json:"img_back_url" validate:"required"`
	Difficulty       string  `json:"difficulty" validate:"required,oneof=easy medium hard"`
	PointsToEarn     int     `json:"points_to_earn" validate:"required"`
	Index            int     `json:"index"`
	IsPublic         bool    `json:"is_public"`
}
