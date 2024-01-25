package types

// Representacion del test module para el frontend

type TestModule struct {
	ID                        uint                       `json:"id"`
	CreatedAt                 string                     `json:"created_at"`
	ModuleID                  uint                       `json:"module_id"`
	Module                    Module                     `json:"module"`
	Started                   *string                    `json:"started"`
	Finished                  *string                    `json:"finished"`
	Qualification             float32                    `json:"qualification"`
	TestModuleQuestionAnswers []TestModuleQuestionAnswer `json:"test_module_question_answers,omitempty"`
}

// pregunta y respuesta del usuario para el test module
type TestModuleQuestionAnswer struct {
	Question   Question   `json:"question"`
	AnswerUser AnswerUser `json:"answer_user"`
}

type AnswerUser struct {
	// TODO: completar
	AnswerUserID uint      `json:"answer_user_id"`
	TestModuleID uint      `json:"test_module_id"`
	QuestionID   uint      `json:"question_id,omitempty"`
	Question     *Question `json:"question,omitempty"`
	AnswerID     uint      `json:"answer_id"`
	Answer       *Answer   `json:"answer"` // Objeto donde realmente esta la respuesta del usuario
	Responded    bool      `json:"responded"`
	Score        int       `json:"score"`
	IsCorrect    bool      `json:"is_correct"`
	Feedback     string    `json:"feedback"`
	ChatIssueID  *uint     `json:"chat_issue_id"`
}

type FinishTest struct {
	Finish        string  `json:"finish"`
	Qualification float32 `json:"qualification"`
	TestID        uint    `json:"test_id"`
}
