package types

import "errors"

type GPT struct {
	Context      string `json:"context"`
	TypeQuestion string `json:"type_question"` // true_false, multi_choice_text, multi_choice_abc, complete_word, order_word
}

func (g *GPT) Validate() error {
	if g.Context == "" {
		return errors.New("context is required")
	}

	if g.TypeQuestion == "" {
		return errors.New("type_question is required")
	}

	if g.TypeQuestion != QuestionTypeTrueOrFalse && g.TypeQuestion != QuestionTypeMultiChoiceText && g.TypeQuestion != "multi_choice_abc" && g.TypeQuestion != "complete_word" && g.TypeQuestion != "order_word" {
		return errors.New("type_question must be one of: true_or_false, multi_choice_text, multi_choice_abc, complete_word, order_word")
	}

	return nil
}

type ResponseIA struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Model        string `json:"model"`
	Server       string `json:"server"`
}

type RequestGPT struct {
	Request string `json:"request"`
}

type ResponseGPT struct {
	Response string `json:"response"`
}
