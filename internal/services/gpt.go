package services

import (
	"Proyectos-UTEQ/api-ortografia/internal/data"
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/spf13/viper"
	"image/png"
	"time"
)

const (
	NAME_BUCKET = "poliword-bucket-2"
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
			Content: "Eres un asistente para estudiante de escuela, donde los estudiantes están aprendiendo ortografía. La respuestas que me debes que dar debe solo tener entre 150 a 250 caracteres.",
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

func (g *ServiceGPT) GenerateQuestion(typeQuestion string, text string, model int) (*types.Question, error) {
	client := openai.NewClient(g.config.GetString("APP_OPENAI_API_KEY"))

	t := openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionDefinition{
			Name:        "generate_question",
			Description: "Genera preguntas para el aprendizaje ortografía",
		},
	}
	switch typeQuestion {
	case types.QuestionTypeTrueOrFalse:
		t.Function.Parameters = jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"text_root": {
					Type:        jsonschema.String,
					Description: "El enunciado de la pregunta, la pregunta debe ser de verdadero o falso.",
				},
				"difficulty": {
					Type:        jsonschema.Integer,
					Description: "El nivel de dificultad de la pregunta, este campo tiene un rango de 1 a 10",
				},
				"answer": {
					Type:        jsonschema.Boolean,
					Description: "La respuesta correcta de la pregunta",
				},
			},
			Required: []string{"text_root", "difficulty", "answer"},
		}
	case types.QuestionTypeMultiChoiceText:
		t.Function.Parameters = jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"text_root": {
					Type:        jsonschema.String,
					Description: "El enunciado de la pregunta, la pregunta debe ser corta y muy entendible por un niño. por ejemplo ¿Las palabras agudas llevan tilde en la última sílaba si terminan en vocal, n o s?",
				},
				"difficulty": {
					Type:        jsonschema.Integer,
					Description: "El nivel de dificultad de la pregunta, este campo tiene un rango de 1 a 10",
				},
				"options": {
					Type:        jsonschema.Array,
					Description: "Las opciones de la preguntas, una de las opciones debe ser la correcta.",
					Items: &jsonschema.Definition{
						Type:        jsonschema.String,
						Description: "La opción de la pregunta",
					},
				},
				"answer": {
					Type:        jsonschema.String,
					Description: "La respuesta correcta de la pregunta",
				},
			},
			Required: []string{"text_root", "difficulty", "options", "answer"},
		}
	case types.QuestionTypeOrderWord:
		t.Function.Parameters = jsonschema.Definition{
			Type:        jsonschema.Object,
			Description: "Genera preguntas de ordenar palabras, para formar una oración. por ejemplo ¿Ordene las palabras, para formar una oración?",
			Properties: map[string]jsonschema.Definition{
				"text_root": {
					Type:        jsonschema.String,
					Description: "El enunciado de la pregunta, la pregunta es de tipo ordenar palabras por tal motivo la pregunta debe ser parecido a esto: ¿Ordene las palabras, para formar una oración?",
				},
				"difficulty": {
					Type:        jsonschema.Integer,
					Description: "El nivel de dificultad de la pregunta, este campo tiene un rango de 1 a 10",
				},
				"options": {
					Type:        jsonschema.Array,
					Description: "Un array de 3 a 5 palabras que deberán formar una oración, estas palabras deben de estar en el orden correcto.",
					Items: &jsonschema.Definition{
						Type:        jsonschema.String,
						Description: "Una de las palabras que se deberá ordenar",
					},
				},
			},
			Required: []string{"text_root", "difficulty", "options"},
		}

	case types.QuestionTypeCompleteWord:
		t.Function.Parameters = jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"text_root": {
					Type:        jsonschema.String,
					Description: "El enunciado de la pregunta, la pregunta es de tipo completar una palabra, por ejemplo: El perro persiguió al __________ por el jardín.",
				},
				"difficulty": {
					Type:        jsonschema.Integer,
					Description: "El nivel de dificultad de la pregunta, este campo tiene un rango de 1 a 10",
				},
				"hind": {
					Type:        jsonschema.String,
					Description: "Pista para que el niño pueda completar la palabra",
				},
				"answer": {
					Type:        jsonschema.String,
					Description: "Respuesta correcta de la pregunta",
				},
			},
			Required: []string{"text_root", "difficulty", "hind", "answer"},
		}
	}

	// creamos un dialogo
	dialogMessage := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Eres un asistente que genera preguntas y respuestas sobre ortografiá",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Te doy un poco de contexto con el siguiente texto: %s", text),
		},
	}
	modelVersion := openai.GPT3Dot5Turbo
	if model == 4 {
		modelVersion = openai.GPT4
	}
	// Iniciamos la comunicación
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelVersion,
			Messages: dialogMessage,
			Tools:    []openai.Tool{t},
		},
	)

	if err != nil || len(resp.Choices) == 0 {
		return nil, err
	}

	if resp.Choices[0].Message.ToolCalls == nil {
		return nil, fmt.Errorf("no se pudo generar la respuesta")
	}

	msg := resp.Choices[0].Message.ToolCalls[0].Function.Arguments
	fmt.Println(msg)

	// transformamos el string a un struct
	// el tipo de la variable depende de la pregunta
	var pregunta types.Questioner
	switch typeQuestion {
	case types.QuestionTypeTrueOrFalse:
		target := &types.QuestionTrueOrFalse{}
		err = json.Unmarshal([]byte(msg), target)
		if err != nil {
			return nil, err
		}
		pregunta = target
	case types.QuestionTypeMultiChoiceText:
		target := &types.QuestionMultiChoiceText{}
		err = json.Unmarshal([]byte(msg), target)
		if err != nil {
			return nil, err
		}
		pregunta = target
	case types.QuestionTypeOrderWord:
		target := &types.QuestionOrderWord{}
		err = json.Unmarshal([]byte(msg), target)
		if err != nil {
			return nil, err
		}
		pregunta = target
	case types.QuestionTypeCompleteWord:
		target := &types.QuestionCompleteWord{}
		err = json.Unmarshal([]byte(msg), target)
		if err != nil {
			return nil, err
		}
		pregunta = target
	}

	return pregunta.ToQuestion(), nil
}

func (g *ServiceGPT) GenerateImage(text string, model int) (string, error) {

	KEY := g.config.GetString("APP_OPENAI_API_KEY")
	c := openai.NewClient(KEY)
	ctx := context.Background()
	modelSelect := ""
	switch model {
	case 2:
		modelSelect = openai.CreateImageModelDallE2
	case 3:
		modelSelect = openai.CreateImageModelDallE3
	default:
		modelSelect = openai.CreateImageModelDallE2
	}

	reqBase64 := openai.ImageRequest{
		Prompt:         text,
		Model:          modelSelect,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
		Quality:        openai.CreateImageQualityStandard,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		return "", err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		return "", err
	}

	// Creamos el cliente para el bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	defer client.Close()

	// construimos el nombre del archivo.
	objectName := fmt.Sprintf("ia_dalle_%d_%s_%s.png", model, time.Now().Format("2006-01-02"), uuid.NewString())

	// Abrir un escritor para google cloud.
	wc := client.Bucket(NAME_BUCKET).Object(objectName).NewWriter(ctx)
	//wc.ContentType = "image/png"
	//wc.Metadata = map[string]string{
	//	"x-goog-meta-test": "data",
	//}

	defer wc.Close()

	// En caso de querer guardar en el file system del backend
	//file, err := os.Create("example.png")
	//if err != nil {
	//	return "", err
	//}
	//defer file.Close()

	if err := png.Encode(wc, imgData); err != nil {
		fmt.Println(err)
		return "", err
	}

	//if _, err := io.Copy(wc, file); err != nil {
	//	fmt.Printf("io.Copy failed: %v\n", err)
	//	return "", err
	//}
	fmt.Println("The image was saved as example.png")

	// construimos la url del objeto
	objectURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", NAME_BUCKET, objectName)
	return objectURL, nil
}
