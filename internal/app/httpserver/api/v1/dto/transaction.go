package dto

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

// swagger:model
type TransactionResponse struct {
	OperationType          string    `json:"operation_type"`
	Amount                 uint64    `json:"amount"`
	SenderWalletId         *string   `json:"sender_wallet_id,omitempty"`
	SenderWalletBalance    *uint64   `json:"sender_wallet_balance,omitempty"`
	SenderWalletMe         bool      `json:"sender_wallet_me,omitempty"`
	RecipientWalletId      *string   `json:"recipient_wallet_id,omitempty"`
	RecipientWalletBalance *uint64   `json:"recipient_wallet_balance,omitempty"`
	RecipientWalletMe      bool      `json:"recipient_wallet_me,omitempty"`
	Balance                uint64    `json:"balance"`
	ProcessedAt            time.Time `json:"processed_at"`
}

type TransactionsResponse []*TransactionResponse

func (resp *TransactionsResponse) ToJson(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(resp)
}

func (resp *TransactionsResponse) ToCsv(writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	header := []string{
		"OperationType", "Amount",
		"SenderWalletId", "SenderWalletBalance", "SenderWalletMe",
		"RecipientWalletId", "RecipientWalletBalance", "RecipientWalletMe",
		"Balance", "ProcessedAt",
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}
	for _, transactionResponse := range *resp {
		var record []string
		record = append(record, transactionResponse.OperationType)
		record = append(record, strconv.FormatUint(transactionResponse.Amount, 10))
		if transactionResponse.SenderWalletId != nil {
			record = append(record, *transactionResponse.SenderWalletId)
		} else {
			record = append(record, "NULL")
		}
		if transactionResponse.SenderWalletBalance != nil {
			record = append(record, strconv.FormatUint(*transactionResponse.SenderWalletBalance, 10))
		} else {
			record = append(record, "NULL")
		}
		if transactionResponse.SenderWalletMe {
			record = append(record, strconv.FormatBool(transactionResponse.SenderWalletMe))
		} else {
			record = append(record, "NULL")
		}
		if transactionResponse.RecipientWalletId != nil {
			record = append(record, *transactionResponse.RecipientWalletId)
		} else {
			record = append(record, "NULL")
		}
		if transactionResponse.RecipientWalletBalance != nil {
			record = append(record, strconv.FormatUint(*transactionResponse.RecipientWalletBalance, 10))
		} else {
			record = append(record, "NULL")
		}
		if transactionResponse.RecipientWalletMe {
			record = append(record, strconv.FormatBool(transactionResponse.RecipientWalletMe))
		} else {
			record = append(record, "NULL")
		}
		record = append(record, strconv.FormatUint(transactionResponse.Balance, 10))
		record = append(record, transactionResponse.ProcessedAt.Format(time.RFC3339Nano))
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}
