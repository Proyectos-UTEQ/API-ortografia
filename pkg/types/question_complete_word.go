package types

type QuestionCompleteWord struct {
	TextRoot   string `json:"text_root"`
	Difficulty int    `json:"difficulty"`
	Hind       string `json:"hind"`
	Answer     string `json:"answer"`
}

func (qCompleteWord *QuestionCompleteWord) ToQuestion() *Question {
	return &Question{
		TextRoot:     qCompleteWord.TextRoot,
		Difficulty:   qCompleteWord.Difficulty,
		TypeQuestion: QuestionTypeCompleteWord,
		Options: Options{
			Hind: qCompleteWord.Hind,
		},
		CorrectAnswer: &Answer{
			TextOptions: []string{qCompleteWord.Answer},
		},
	}
}
