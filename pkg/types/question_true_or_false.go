package types

type QuestionTrueOrFalse struct {
	TextRoot   string `json:"text_root"`
	Difficulty int    `json:"difficulty"`
	Answer     bool   `json:"answer"`
}

// ToQuestion Convierte una pregunta de TrueOrFalse a Question.
func (qTrueOrFalse *QuestionTrueOrFalse) ToQuestion() *Question {
	return &Question{
		TextRoot:     qTrueOrFalse.TextRoot,
		Difficulty:   qTrueOrFalse.Difficulty,
		TypeQuestion: QuestionTypeTrueOrFalse,
		CorrectAnswer: &Answer{
			TrueOrFalse: qTrueOrFalse.Answer,
		},
	}
}
