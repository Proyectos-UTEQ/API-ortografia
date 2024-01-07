package types

import "errors"

// Datos para que un usuario se puede suscribir a un modulo.
type ReqSubscription struct {
	Code string `json:"code"`
}

func (r *ReqSubscription) Validate() error {
	if r.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
