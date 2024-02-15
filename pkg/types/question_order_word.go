package types

import (
	"math/rand"
	"time"
)

type QuestionOrderWord struct {
	TextRoot   string   `json:"text_root"`
	Difficulty int      `json:"difficulty"`
	Options    []string `json:"options"`
}

func (qOrderWord *QuestionOrderWord) ToQuestion() *Question {
	// desordenar las opciones
	rand.Seed(time.Now().UnixNano())
	shuffledArray := make([]string, len(qOrderWord.Options))
	copy(shuffledArray, qOrderWord.Options)
	rand.Shuffle(len(shuffledArray), func(i, j int) {
		shuffledArray[i], shuffledArray[j] = shuffledArray[j], shuffledArray[i]
	})
	return &Question{
		TextRoot:     qOrderWord.TextRoot,
		Difficulty:   qOrderWord.Difficulty,
		TypeQuestion: QuestionTypeOrderWord,
		Options: Options{
			TextOptions: qOrderWord.Options,
		},
		CorrectAnswer: &Answer{
			TextOptions: qOrderWord.Options,
		},
	}
}
