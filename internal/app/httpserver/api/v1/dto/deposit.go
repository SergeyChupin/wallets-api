package dto

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type DepositRequest struct {
	Amount uint64 `json:"amount" validate:"required"`
}

func (req *DepositRequest) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(req)
}

func (req *DepositRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}

// swagger:model
type DepositResponse struct {
	Balance uint64 `json:"balance"`
}

func (resp *DepositResponse) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(resp)
}
