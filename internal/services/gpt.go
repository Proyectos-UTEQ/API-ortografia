package services

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/spf13/viper"
)

type ServiceGPT struct {
	config *viper.Viper
}

func NewGPT(config *viper.Viper) *ServiceGPT {
	return &ServiceGPT{
		config: config,
	}
}

func (g *ServiceGPT) GenerateResponse(msg string) (string, error) {
	client := openai.NewClient(g.config.GetString("APP_OPENAI_API_KEY"))
	rest, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Eres un asistente para estudiante de escuela, donde los estudiantes están aprendiendo ortografía.",
				},
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

func (g *ServiceGPT) GenerateFeedbackForQuestion(answerUser *data.AnswerUser) error {
	client := openai.NewClient(g.config.GetString("APP_OPENAI_API_KEY"))

	dialogMessage := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Eres un asistente para estudiante de escuela, donde los estudiantes están aprendiendo ortografía. La respuestas que me debes que dar debe solo tener entre 100 a 150 caracteres.",
		},
	}

	switch answerUser.Question.TypeQuestion {
	case "true_or_false":
		dialogMessage = append(dialogMessage, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Necesito que me des retroalimentación para la pregunta, respuesta correcta y la respuesta del estudiante, a continuación te dejo los datos. Pregunta: %s. Respuesta correcta: %t. Respuesta del estudiante: %t", answerUser.Question.TextRoot, answerUser.Question.CorrectAnswer.TrueOrFalse, answerUser.Answer.TrueOrFalse),
		})
	case "multi_choice_text":
		dialogMessage = append(dialogMessage, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Necesito que me des retroalimentación para la pregunta, respuesta correcta y la respuesta del estudiante, a continuación te dejo los datos. Pregunta: %s. Respuesta correcta: %s. Respuesta del estudiante: %s", answerUser.Question.TextRoot, answerUser.Question.CorrectAnswer.TextOptions, answerUser.Answer.TextOptions),
		})
	case "complete_word":
		dialogMessage = append(dialogMessage, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Necesito que me des retroalimentación para la pregunta de completación, respuesta correcta y la respuesta del estudiante, a continuación te dejo los datos. Pregunta: %s. Respuesta correcta: %s. Respuesta del estudiante: %s", answerUser.Question.TextRoot, answerUser.Question.CorrectAnswer.TextToComplete, answerUser.Answer.TextToComplete),
		})
	case "order_word":
		dialogMessage = append(dialogMessage, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Necesito que me des retroalimentación para la pregunta de orden de palabras, respuesta correcta y la respuesta del estudiante, a continuación te dejo los datos. Pregunta: %s. Respuesta correcta: %s. Respuesta del estudiante: %s", answerUser.Question.TextRoot, answerUser.Question.CorrectAnswer.TextOptions, answerUser.Answer.TextOptions),
		})
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: dialogMessage,
		},
	)
	if err != nil {
		return err
	}

	answerUser.Feedback = resp.Choices[0].Message.Content
	return nil
}

func (g *ServiceGPT) GenerateQuestion(typeQuestion string, text string) (*types.Question, error) {
	client := openai.NewClient(g.config.GetString("APP_OPENAI_API_KEY"))

	t := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionDefinition{
			Name:        "generate_question",
			Description: "Genera una pregunta de tipo true_false, multi_choice_text, multi_choice_abc, complete_word, order_word",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"text_root": {
						Type:        jsonschema.String,
						Description: "El enunciado de la pregunta, la pregunta debe ser corta y muy entendible por un niño. por ejemplo ¿En las palabras esdrújula donde llevan tilde?",
					},
					"difficulty": {
						Type:        jsonschema.Integer,
						Description: "El nivel de dificultad de la pregunta, este campo tiene un rango de 1 a 10",
					},
					"type_question": {
						Type:        jsonschema.String,
						Description: "El tipo de pregunta, puede ser true_false, multi_choice_text, multi_choice_abc, complete_word, order_word",
					},
					"options": {
						Type:        jsonschema.Object,
						Description: "Las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word",
						Properties: map[string]jsonschema.Definition{
							"select_mode": {
								Type:        jsonschema.String,
								Enum:        []string{"single", "multiple"},
								Description: "El modo de seleccionar las opciones, solo se puede elegir entre single o multiple, en caso de seleccionar single en el campo correct_answer.text_options solo debe ir una respuesta, en caso de seleccionar multiple en el campo correct_answer.text_options debe ir una lista de respuestas",
							},
							"text_options": {
								Type: jsonschema.Array,
								Items: &jsonschema.Definition{
									Type:        jsonschema.String,
									Description: "Las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word",
								},
							},
							"text_to_complete": {
								Type:        jsonschema.String,
								Description: "El texto que se va a completar, este campo solo se llena en caso de que el tipo de pregunta sea complete_word",
							},
							"hind": {
								Type:        jsonschema.String,
								Description: "La pista para resolver la pregunta",
							},
						},
					},
					"correct_answer": {
						Type:        jsonschema.Object,
						Description: "La respuesta correcta para la pregunta que esta en el enunciado, este campo solo se llena en caso de que el tipo de pregunta sea true_false, multi_choice_text, multi_choice_abc, complete_word, order_word",
						Properties: map[string]jsonschema.Definition{
							"true_or_false": {
								Type:        jsonschema.Boolean,
								Description: "Si la respuesta es correcta o no, en caso de ser type_question true_false",
							},
							"text_options": {
								Type: jsonschema.Array,
								Items: &jsonschema.Definition{
									Type:        jsonschema.String,
									Description: "Las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word",
								},
							},
							"text_to_complete": {
								Type: jsonschema.Array,
								Items: &jsonschema.Definition{
									Type:        jsonschema.String,
									Description: "Las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea text_to_complete",
								},
							},
						},
					},
				},
				Required: []string{"text_root", "difficulty", "type_question", "options", "correct_answer"},
			},
		},
	}

	// creamos un dialogo
	dialogMessage := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Generame una pregunta de tipo %s, las preguntas las debes sacar del siguiente texto: %s", typeQuestion, text)},
	}

	// Iniciamos la comunicación
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: dialogMessage,
			Tools:    []openai.Tool{t},
		},
	)

	if err != nil || len(resp.Choices) == 0 {
		return nil, err
	}

	msg := resp.Choices[0].Message.ToolCalls[0].Function.Arguments

	pregunta := &types.Question{}
	err = json.Unmarshal([]byte(msg), pregunta)
	if err != nil {
		return nil, err
	}

	return pregunta, nil
}

func (g *ServiceGPT) GenerateImage(text string) error {

	return nil
}
