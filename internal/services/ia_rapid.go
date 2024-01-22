package services

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

const SchemeJson = `{
    "module_id": 2,
    "text_root": "enunciado de la pregunta",
    "difficulty": 1,
    "type_question": "order_word",
    "question_answer": {
        "select_mode": "",
        "text_options": ["palabra 2", "palabra 1", ...],
        "text_to_complete": "",
        "hind": "Pista para solucionar la pregunta"
    },
    "correct_answer": {
        "true_or_false": true,
        "text_opcions": ["palabra 1", "palabra 2"],
        "text_to_complete": []
    }
}`

type ServicoIA struct {
	config *viper.Viper
}

func NewServicoIA(config *viper.Viper) *ServicoIA {
	return &ServicoIA{
		config: config,
	}
}

// estructura para enviar la petición a la IA
type ReqIA struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func (s *ServicoIA) IARapidTest(msg string) (*types.Question, error) {
	reqIA := []ReqIA{
		{
			Role:    "system",
			Content: "Tu eres un profesor universitario que genera preguntas muy bien pensadas",
		},
		{
			Role: "system",
			Content: `Te explicare un poco el esquema que te estoy pasando, el campo module_id simplemente es el id del modulo al cual se agregar esta pregunta lo puedes dejar en 0,
			el campo difficulty es el nivel de dificulta que tiene la pregunta, este campo tiene un rago de 1 a 10, 
			en el campo type_question puedes elegir entre: true_false, multi_choice_text, multi_choice_abc, complete_word, order_word.
			en el campo question_answer es el objeto que contiene la informacion de la respuesta, en el campo select_mode puedes elegir entre: single, multiple.
			en el campo text_options es un array de string el cual tendra las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word
			en el campo text_to_complete es un string el cual tendra una oración donde se debe completar con pabras.
			el campo hind es de tipo string el cual tendra una pista para solucionar la pregunta.
			pasando al campo correct_answer es el objeto que contiene la informacion de la respuesta correcta, en el campo true_or_false puedes elegir entre: true, false.
			en el campo text_opcions es un array de string el cual tendra las palabras correctas, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text, multi_choice_abc y order_word,
			en caso de ser order_word el array de string tendra el orden correcto de las palabras.
			en el campo text_to_complete es un array de string el cual tendra las palabras que se deben completar con pabras.`,
		},
		{
			Role:    "system",
			Content: "Proporcione su respuesta estrictamente como un objeto JSON con el siguiente esquema: " + SchemeJson,
		},
		{
			Role:    "user",
			Content: "generame un pregunta, basandote en el siguiente enunciado: " + msg,
		},
	}

	// convertimos el objecto.
	reqIAbyte, err := json.Marshal(reqIA)
	if err != nil {
		return nil, err
	}

	// contruimos el payload.
	payload := strings.NewReader(string(reqIAbyte))

	url := s.config.GetString("APP_IA_RAPID_URL")

	req, _ := http.NewRequest("POST", url, payload)

	// agregamos los headers.
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Key", s.config.GetString("APP_RAPIDAPI_KEY"))
	req.Header.Add("X-RapidAPI-Host", "chatgpt-api8.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return nil, errors.New("error en la petición")
	}

	responseIA := &types.ResponseIA{}

	err = json.Unmarshal(body, responseIA)
	if err != nil {
		return nil, err
	}

	fmt.Println(responseIA)
	// Realizar un tratamiento a la respuesta

	startIndex := strings.Index(responseIA.Text, "{")
	if startIndex == -1 {
		return nil, errors.New("error en la respuesta")
	}
	endIndex := strings.LastIndex(responseIA.Text, "}") + 1

	jsonString := responseIA.Text[startIndex:endIndex]

	pregunta := &types.Question{}
	err = json.Unmarshal([]byte(jsonString), pregunta)
	if err != nil {
		fmt.Println()
		return nil, err
	}

	fmt.Println(pregunta)

	return pregunta, nil
}
