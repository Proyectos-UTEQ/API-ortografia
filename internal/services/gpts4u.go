package services

import (
	"Proyectos-UTEQ/api-ortografia/pkg/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type ServiceNewGpts4u struct {
	config *viper.Viper
}

func NewServiceNewGpts4u(config *viper.Viper) *ServiceNewGpts4u {
	return &ServiceNewGpts4u{
		config: config,
	}
}

type ReqNewChatGPT struct {
	Statu    string `json:"status"`
	Chattext string `json:"chattext"`
	Action   string `json:"action"`
}

func (s *ServiceNewGpts4u) GenerateQuestion(typequestion string, text string) (*types.Question, error) {

	reqIA := []ReqIA{
		{
			Role:    "system",
			Content: fmt.Sprintf("Generame una pregunta de tipo %s, basandote en el siguiente enunciado: %s, necesito que las respuesta sea en formato json como el siguiente: %s, en el campo type_question puedes elegir entre: true_false, multi_choice_text, multi_choice_abc, complete_word, order_word, en el campo correct_answer es el objeto que contiene la informacion de la respuesta dependiendo del type_question se llenan los campos true_or_false, text_opcions, text_to_complete", typequestion, text, SchemeJson),
		},
	}

	reqbyte, err := json.Marshal(reqIA)
	if err != nil {
		return nil, err
	}

	payload := strings.NewReader(string(reqbyte))

	// Create a new request using http
	queryParams := url.Values{}
	queryParams.Add("role", "user")
	queryParams.Add("content", fmt.Sprintf(`Generame una pregunta de tipo %s con su respentiva respuesta, basandote en el siguiente enunciado: %s, necesito que las respuesta sea en formato json como el siguiente: %s
	Te explicare un poco el esquema que te estoy pasando, el campo module_id simplemente es el id del modulo al cual se agregar esta pregunta lo puedes dejar en 0,
	el campo difficulty es el nivel de dificulta que tiene la pregunta, este campo tiene un rago de 1 a 10, 
	en el campo type_question puedes elegir entre: true_false, multi_choice_text, multi_choice_abc, complete_word, order_word.
	en el campo question_answer es el objeto que contiene la informacion de la respuesta, en el campo select_mode puedes elegir entre: single, multiple.
	en el campo text_options es un array de string el cual tendra las opciones de la respuesta, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text o multi_choice_abc, order_word
	en el campo text_to_complete es un string el cual tendra una oración donde se debe completar con pabras.
	el campo hind es de tipo string el cual tendra una pista para solucionar la pregunta.
	pasando al campo correct_answer es el objeto que contiene la informacion de la respuesta correcta, 
	en el campo correct_answer.true_or_false puedes elegir entre: true, false. 
	en el campo text_opcions es un array de string el cual tendra las palabras correctas, este campo solo se llena en caso de que el tipo de pregunta sea multi_choice_text, multi_choice_abc y order_word,
	en caso de ser order_word el array de string tendra el orden correcto de las palabras.
	en el campo text_to_complete es un array de string el cual tendra las palabras que se deben completar con pabras.
	.
	`, typequestion, text, SchemeJson))

	url := "https://gpts4u.p.rapidapi.com/bingChat" + "?" + queryParams.Encode()
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	// Add the headers
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("X-RapidAPI-Key", s.config.GetString("APP_RAPIDAPI_KEY"))
	req.Header.Add("X-RapidAPI-Host", "gpts4u.p.rapidapi.com")

	fmt.Println("Realizando peticion")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("Respuesta recibida")

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, err
	}

	jsonString, err := getJson(string(body))
	if err != nil {
		return nil, err
	}

	fmt.Println(jsonString)

	jsonString = strings.ReplaceAll(jsonString, "\\n", "")
	jsonString = strings.ReplaceAll(jsonString, "\\", "")

	fmt.Println(jsonString)

	pregunta := &types.Question{}

	err = json.Unmarshal([]byte(jsonString), pregunta)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pregunta, nil

}

func getJson(text string) (string, error) {

	startIndex := strings.Index(text, "{")
	if startIndex == -1 {
		return "", errors.New("error en la respuesta")
	}

	endIndex := strings.LastIndex(text, "}") + 1

	jsonString := text[startIndex:endIndex]
	return jsonString, nil
}
