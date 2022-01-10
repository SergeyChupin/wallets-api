package dto

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

// swagger:model
type TransferRequest struct {
	Amount         uint64 `json:"amount" validate:"required,gt=0"`
	SenderWalletId string `json:"sender_wallet_id" validate:"required"`
}

func (req *TransferRequest) FromJson(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(req)
}

func (req *TransferRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}

// swagger:model
type TransferResponse struct {
	SenderWalletBalance uint64 `json:"sender_wallet_balance"`
	Balance             uint64 `json:"balance"`
}

func (resp *TransferResponse) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(resp)
}
