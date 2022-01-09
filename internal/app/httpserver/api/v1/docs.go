// Wallets API
//
// Documentation for Wallets API
//
//	Schemes: http
//	BasePath: /api/v1
//	Version: 1.0.0
//
// swagger:meta
package v1

import (
	"time"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/api/v1/dto"
)

// swagger:parameters createWallet
type createWalletRequest struct {
	// in: body
	Body dto.CreateWalletRequest `json:"body"`
}

// swagger:response createWalletResponse
type createWalletResponse struct {
	// in: body
	Body dto.CreateWalletResponse `json:"body"`
}

// swagger:parameters deposit
type depositRequest struct {
	// in: body
	Body dto.DepositRequest `json:"body"`
}

// swagger:response depositResponse
type depositResponse struct {
	// in: body
	Body dto.DepositResponse `json:"body"`
}

// swagger:parameters transfer
type transferRequest struct {
	// in: body
	Body dto.TransferRequest `json:"body"`
}

// swagger:response transferResponse
type transferResponse struct {
	// in: body
	Body dto.TransferResponse `json:"body"`
}

// swagger:parameters deposit transfer getTransactions
type walletID struct {
	// in: path
	ID string `json:"id"`
}

// swagger:parameters getTransactions
type getTransactions struct {
	// in: query
	Limit int `json:"limit"`
	// in: query
	Offset int `json:"offset"`
	// in: query
	OperationType string `json:"operation_type"`
	// in: query
	ProcessedAtGte time.Time `json:"processed_at.gte"`
	// in: query
	ProcessedAtLte time.Time `json:"processed_at.lte"`
}

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body dto.ErrorResponse
}

// swagger:response transactionsResponse
type transactionsResponse struct {
	// in: body
	Body []dto.TransactionResponse
}
