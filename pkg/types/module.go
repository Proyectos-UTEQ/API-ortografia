package types

// Representacion de un modulo para el frontend
type Module struct {
	ID               uint    `json:"id"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	CreateBy         UserAPI `json:"create_by"`
	Title            string  `json:"title"`
	ShortDescription string  `json:"short_description"`
	TextRoot         string  `json:"text_root"`
	ImgBackURL       string  `json:"img_back_url"`
	Difficulty       string  `json:"difficulty"`
	PointsToEarn     string  `json:"points_to_earn"`
	Index            int     `json:"index"`
	IsPublic         bool    `json:"is_public"`
}
