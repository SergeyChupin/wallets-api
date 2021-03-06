package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/api/v1/dto"
	"github.com/SergeyChupin/wallets-api/internal/model"
	"github.com/SergeyChupin/wallets-api/internal/service"
	"github.com/gorilla/mux"
)

type walletsApi struct {
	logger        *log.Logger
	walletService service.WalletService
}

func NewWalletsApi(logger *log.Logger, router *mux.Router, walletService service.WalletService) {
	walletsApi := &walletsApi{
		logger:        logger,
		walletService: walletService,
	}
	router.HandleFunc("/wallets", walletsApi.CreateWallet).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/deposit", walletsApi.Deposit).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transfer", walletsApi.Transfer).Methods(http.MethodPost)
	router.HandleFunc("/wallets/{id}/transactions", walletsApi.GetTransactions).Methods(http.MethodGet)
}

// swagger:route POST /wallets WalletsAPI createWallet
// Create new wallet
//
// consumes:
//	- application/json
// produces:
// 	- application/json
//
// responses:
//	200: createWalletResponse
//  400: errorResponse
//  500: errorResponse
func (walletsApi *walletsApi) CreateWallet(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var reqData dto.CreateWalletRequest
	err := reqData.FromJson(req.Body)
	if err != nil {
		walletsApi.logger.Println("walletsApi - CreateWallet - reqData.FromJson:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	if err = reqData.Validate(); err != nil {
		walletsApi.logger.Println("walletsApi - CreateWallet - reqData.Validate:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := walletsApi.walletService.CreateWallet(
		model.Wallet{
			Name:     reqData.Name,
			Currency: reqData.Currency,
		},
	)
	if err != nil {
		walletsApi.logger.Println("walletsApi - CreateWallet - walletsApi.walletService.CreateWallet:", err)
		writeError(rw, "unable to create wallet", http.StatusInternalServerError)
		return
	}
	respData := dto.CreateWalletResponse{ID: id}
	if err = respData.ToJson(rw); err != nil {
		walletsApi.logger.Println("walletsApi - CreateWallet - respData.ToJson:", err)
		writeError(rw, "internal error", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /wallets/{id}/deposit WalletsAPI deposit
// Deposit money to wallet
//
// consumes:
//	- application/json
// produces:
// 	- application/json
//
// responses:
//	200: depositResponse
//  400: errorResponse
//  500: errorResponse
func (walletsApi *walletsApi) Deposit(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	id := getWalletId(req)
	var reqData dto.DepositRequest
	if err := reqData.FromJson(req.Body); err != nil {
		walletsApi.logger.Println("walletsApi - Deposit - reqData.FromJson:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := reqData.Validate(); err != nil {
		walletsApi.logger.Println("walletsApi - Deposit - reqData.Validate:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	depositTransaction, err := walletsApi.walletService.Deposit(
		id, reqData.Amount,
	)
	if err != nil {
		walletsApi.logger.Println("walletsApi - Deposit - walletsApi.walletService.Deposit:", err)
		writeError(rw, "unable to deposit wallet", http.StatusInternalServerError)
		return
	}
	respData := dto.DepositResponse{Balance: depositTransaction.RecipientWallet.Balance}
	if err = respData.ToJson(rw); err != nil {
		walletsApi.logger.Println("walletsApi - Deposit - respData.ToJson:", err)
		writeError(rw, "internal error", http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /wallets/{id}/transfer WalletsAPI transfer
// Transfer money between wallets
//
// consumes:
//	- application/json
// produces:
// 	- application/json
//
// responses:
//	200: transferResponse
//  400: errorResponse
//  500: errorResponse
func (walletsApi *walletsApi) Transfer(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	id := getWalletId(req)
	var reqData dto.TransferRequest
	if err := reqData.FromJson(req.Body); err != nil {
		walletsApi.logger.Println("walletsApi - Transfer - reqData.FromJson:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := reqData.Validate(); err != nil {
		walletsApi.logger.Println("walletsApi - Transfer - reqData.Validate:", err)
		writeError(rw, "invalid request body", http.StatusBadRequest)
		return
	}
	transferTransaction, err := walletsApi.walletService.Transfer(
		reqData.SenderWalletId, id, reqData.Amount,
	)
	if err != nil {
		walletsApi.logger.Println("walletsApi - Transfer - walletsApi.walletService.Transfer:", err)
		writeError(rw, "unable to transfer money between wallets", http.StatusInternalServerError)
		return
	}
	respData := dto.TransferResponse{
		SenderWalletBalance: transferTransaction.SenderWallet.Balance,
		Balance:             transferTransaction.RecipientWallet.Balance,
	}
	if err = respData.ToJson(rw); err != nil {
		walletsApi.logger.Println("walletsApi - Transfer - respData.ToJson:", err)
		writeError(rw, "internal error", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /wallets/{id}/transactions WalletsAPI getTransactions
// Return a list of transactions
//
// produces:
// 	- application/json
//	- text/csv
//
// responses:
//	200: transactionsResponse
//  400: errorResponse
//  406: errorResponse
//  500: errorResponse
func (walletsApi *walletsApi) GetTransactions(rw http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Accept")
	if contentType == "" || contentType == "*/*" {
		contentType = "application/json"
	}
	if contentType != "application/json" && contentType != "text/csv" {
		walletsApi.logger.Println("walletsApi - GetTransactions - invalid header 'Accept'")
		writeError(rw, "invalid header 'Accept'", http.StatusNotAcceptable)
		return
	}
	rw.Header().Set("Content-Type", contentType)
	id := getWalletId(req)
	limit, offset, err := getPagination(req)
	if err != nil {
		walletsApi.logger.Println("walletsApi - GetTransactions - getPagination:", err)
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	filter := model.TransactionFilter{
		WalletId: id,
	}
	operationType := req.URL.Query().Get("operation_type")
	if operationType != "" {
		operationType, err := model.FromString(operationType)
		if err != nil {
			walletsApi.logger.Println("walletsApi - GetTransactions - model.FromString:", err)
			writeError(rw, "invalid query parameter operation_type", http.StatusBadRequest)
			return
		}
		filter.OperationType = operationType
	}
	processedAtGte := req.URL.Query().Get("processed_at.gte")
	if processedAtGte != "" {
		processedAtGte, err := time.Parse(time.RFC3339Nano, processedAtGte)
		if err != nil {
			walletsApi.logger.Println("walletsApi - GetTransactions - time.Parse:", err)
			writeError(rw, "invalid query parameter processed_at.gte", http.StatusBadRequest)
			return
		}
		filter.ProcessedAtGte = processedAtGte
	}
	processedAtLte := req.URL.Query().Get("processed_at.lte")
	if processedAtLte != "" {
		processedAtLte, err := time.Parse(time.RFC3339Nano, processedAtLte)
		if err != nil {
			walletsApi.logger.Println("walletsApi - GetTransactions - time.Parse:", err)
			writeError(rw, "invalid query parameter processed_at.lte", http.StatusBadRequest)
			return
		}
		filter.ProcessedAtLte = processedAtLte
	}
	if !filter.ProcessedAtLte.IsZero() && !filter.ProcessedAtGte.Before(filter.ProcessedAtLte) {
		walletsApi.logger.Println("walletsApi - GetTransactions - invalid time range processed_at")
		writeError(rw, "invalid time range processed_at", http.StatusBadRequest)
		return
	}
	transactions, err := walletsApi.walletService.GetTransactions(
		limit, offset, filter,
	)
	if err != nil {
		walletsApi.logger.Println("walletsApi - GetTransactions - walletsApi.walletService.GetTransactions:", err)
		writeError(rw, "unable to get transactions", http.StatusInternalServerError)
		return
	}
	var respData dto.TransactionsResponse = make([]*dto.TransactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		respItem := &dto.TransactionResponse{
			OperationType: transaction.OperationType.String(),
			Amount:        transaction.Amount,
			Balance:       transaction.RecipientWallet.Balance,
			ProcessedAt:   transaction.ProcessedAt,
		}
		if transaction.SenderWallet != nil {
			if transaction.SenderWallet.ID == id {
				respItem.Balance = transaction.SenderWallet.Balance
				respItem.SenderWalletMe = true
				respItem.RecipientWalletId = &transaction.RecipientWallet.ID
				respItem.RecipientWalletBalance = &transaction.RecipientWallet.Balance
			} else {
				respItem.Balance = transaction.RecipientWallet.Balance
				respItem.RecipientWalletMe = true
				respItem.SenderWalletId = &transaction.SenderWallet.ID
				respItem.SenderWalletBalance = &transaction.SenderWallet.Balance
			}
		}
		respData = append(respData, respItem)
	}
	if contentType == "application/json" {
		if err = respData.ToJson(rw); err != nil {
			walletsApi.logger.Println("walletsApi - GetTransactions - respData.ToJson:", err)
			writeError(rw, "internal error", http.StatusInternalServerError)
			return
		}
	}
	if contentType == "text/csv" {
		if err = respData.ToCsv(rw); err != nil {
			walletsApi.logger.Println("walletsApi - GetTransactions - respData.ToCsv:", err)
			writeError(rw, "internal error", http.StatusInternalServerError)
			return
		}
	}
}

func writeError(rw http.ResponseWriter, message string, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	errRespData := dto.ErrorResponse{Message: message}
	_ = errRespData.ToJson(rw)
}

func getWalletId(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["id"]
}

func getPagination(req *http.Request) (pageLimit int, pageOffset int, err error) {
	pageLimit, pageOffset = -1, -1
	limit := req.URL.Query().Get("limit")
	if limit != "" {
		pageLimit, err = strconv.Atoi(limit)
		if err != nil || pageLimit < -1 {
			err = errors.New("invalid limit query parameter")
			return
		}
	}
	offset := req.URL.Query().Get("offset")
	if offset != "" {
		pageOffset, err = strconv.Atoi(offset)
		if err != nil || pageOffset < -1 {
			err = errors.New("invalid offset query parameter")
			return
		}
	}
	return
}
