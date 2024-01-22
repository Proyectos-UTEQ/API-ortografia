package types

type GPT struct {
	Context string `json:"context"`
}

type ResponseIA struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Model        string `json:"model"`
	Server       string `json:"server"`
}
