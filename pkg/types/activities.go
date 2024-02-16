package types

// Activities se lista por el módulo selection para mostrar las preguntas
// de ese módulo.
type Activities struct {
	ID           uint    `json:"id"` // ID de la pregunta
	TypeQuestion string  `json:"type_question"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    *string `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
	Difficulty   int     `json:"difficulty"`
	TextRoot     string  `json:"text_root"`
}

//*Nombre del tipo de actividad (type_question)
//*Creada por
//*Dificultad
//*Fecha de creación
//*Fecha de actualización
