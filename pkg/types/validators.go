package types

import "github.com/go-playground/validator/v10"

type ErrorField struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func Validate(data interface{}) ([]ErrorField, error) {
	validate := validator.New()
	err := validate.Struct(data)

	if err != nil {

		resp := make([]ErrorField, 0)
		for _, err := range err.(validator.ValidationErrors) {
			resp = append(resp, ErrorField{
				Field: err.Field(),
				Tag:   err.Tag(),
			})

			return resp, err
		}

		return nil, err
	}

	return nil, nil
}
