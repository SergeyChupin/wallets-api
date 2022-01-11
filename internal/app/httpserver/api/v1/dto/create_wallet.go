package dto

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type CreateWalletRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD"`
}

func (req *CreateWalletRequest) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(req)
}

func (req *CreateWalletRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}

// swagger:model
type CreateWalletResponse struct {
	ID string `json:"id"`
}

func (resp *CreateWalletResponse) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(resp)
}
