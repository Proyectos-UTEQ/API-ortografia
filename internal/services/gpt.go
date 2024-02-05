package services

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"

	openai "github.com/sashabaranov/go-openai"
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

func (g *ServiceGPT) GenerateQuestion(typeQuestion string, text string) (*types.Question, error) {
	client := openai.NewClient(g.config.GetString("APP_OPENAI_API_KEY"))
	rest, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo1106,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Eres un asistente especializado en ortografía para generar preguntas estrictamente en formato JSON",
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Te explicare un poco el esquema que te estoy pasando, el campo module_id simplemente es el id del modulo al cual se agregará esta pregunta lo puedes dejar en 0,
							el campo difficulty es el nivel de dificulta que tiene la pregunta, este campo tiene un rango de 1 a 10, 
							en el campo type_question puedes elegir entre: true_false, multi_choice_text, multi_choice_abc, complete_word, order_word.
							en el campo question_answer es el objeto que contiene la información de la respuesta, en el campo select_mode puedes elegir entre: single, multiple.
							en el campo text_options es un array de string el cual tendrá las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word
							en el campo text_to_complete es un string el cual tendrá una oración donde se debe completar con palabras.
							el campo hind es de tipo string el cual tendrá una pista para solucionar la pregunta.
							pasando al campo correct_answer es el objeto que contiene la información de la respuesta correcta, en el campo true_or_false puedes elegir entre: true, false.
							en el campo text_opcions es un array de string el cual tendrá las palabras correctas, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text, multi_choice_abc y order_word,
							en caso de ser order_word el array de string tendrá el orden correcto de las palabras.
							en el campo text_to_complete es un array de string el cual tendrá las palabras que se deben completar con palabras.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Generame una pregunta de tipo %s, basandote en el siguiente enunciado: %s, necesito que las respuesta sea estrictamente en formato json como el siguiente: %s", typeQuestion, text, SchemeJson),
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	question := rest.Choices[0].Message.Content

	startIndex := strings.Index(question, "{")
	if startIndex == -1 {
		return nil, errors.New("error en la respuesta")
	}
	endIndex := strings.LastIndex(question, "}") + 1

	jsonString := question[startIndex:endIndex]

	fmt.Println(jsonString)

	pregunta := &types.Question{}
	err = json.Unmarshal([]byte(jsonString), pregunta)
	if err != nil {
		return nil, err
	}

	return pregunta, nil
}
