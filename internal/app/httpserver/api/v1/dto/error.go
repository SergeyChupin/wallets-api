package dto

import (
	"encoding/json"
	"io"
)

// swagger:model
type ErrorResponse struct {
	Message string `json:"message"`
}

func (errorResponse *ErrorResponse) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(errorResponse)
}
