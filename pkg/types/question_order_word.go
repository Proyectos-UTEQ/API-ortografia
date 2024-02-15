package types

type QuestionOrderWord struct {
	TextRoot   string   `json:"text_root"`
	Difficulty int      `json:"difficulty"`
	Options    []string `json:"options"`
	Answer     []string `json:"answer"`
}

func (qOrderWord *QuestionOrderWord) ToQuestion() *Question {
	return &Question{
		TextRoot:     qOrderWord.TextRoot,
		Difficulty:   qOrderWord.Difficulty,
		TypeQuestion: QuestionTypeOrderWord,
		Options: Options{
			TextOptions: qOrderWord.Options,
		},
		CorrectAnswer: &Answer{
			TextOptions: qOrderWord.Answer,
		},
	}
}
