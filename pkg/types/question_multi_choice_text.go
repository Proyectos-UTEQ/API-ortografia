package types

type QuestionMultiChoiceText struct {
	TextRoot   string   `json:"text_root"`
	Difficulty int      `json:"difficulty"`
	Options    []string `json:"options"`
	Answer     string   `json:"answer"`
}

func (qMultiChoice *QuestionMultiChoiceText) ToQuestion() *Question {
	return &Question{
		TextRoot:     qMultiChoice.TextRoot,
		Difficulty:   qMultiChoice.Difficulty,
		TypeQuestion: QuestionTypeMultiChoiceText,
		Options: Options{
			SelectMode:  "single",
			TextOptions: qMultiChoice.Options,
		},
		CorrectAnswer: &Answer{
			TextOptions: []string{qMultiChoice.Answer},
		},
	}
}
