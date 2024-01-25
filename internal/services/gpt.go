package services

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type ServiceGPT struct {
	Key string
}

func NewGPT(key string) *ServiceGPT {
	return &ServiceGPT{
		Key: key,
	}
}

func (g *ServiceGPT) GPTTest(msg string) (string, error) {
	client := openai.NewClient(g.Key)
	rest, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	msgGPT := rest.Choices[0].Message.Content

	return msgGPT, nil
}
