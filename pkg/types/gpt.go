package types

import "errors"

type GPT struct {
	Context      string `json:"context"`
	TypeQuestion string `json:"type_question"` // true_false, multi_choice_text, multi_choice_abc, complete_word, order_word
	Server       string `json:"server"`        // chatgptapi or gpts4u
}

func (g *GPT) Validate() error {
	if g.Context == "" {
		return errors.New("context is required")
	}

	if g.TypeQuestion == "" {
		return errors.New("type_question is required")
	}

	if g.TypeQuestion != "true_false" && g.TypeQuestion != "multi_choice_text" && g.TypeQuestion != "multi_choice_abc" && g.TypeQuestion != "complete_word" && g.TypeQuestion != "order_word" {
		return errors.New("type_question must be one of: true_false, multi_choice_text, multi_choice_abc, complete_word, order_word")
	}

	if g.Server == "" {
		return errors.New("server is required")
	}

	if g.Server != "chatgptapi" && g.Server != "gpts4u" {
		return errors.New("server must be one of: chatgptapi, newchatgpt")
	}

	return nil
}

type ResponseIA struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Model        string `json:"model"`
	Server       string `json:"server"`
}
