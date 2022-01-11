package v1

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/SergeyChupin/wallets-api/internal/app/httpserver/api/v1/dto"
	"github.com/SergeyChupin/wallets-api/internal/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	logger = log.New(os.Stdout, "wallets-api-testing ", log.LstdFlags)
)

type walletServiceMock struct {
	mock.Mock
}

func (walletService *walletServiceMock) CreateWallet(wallet model.Wallet) (string, error) {
	args := walletService.Called(wallet)
	return args.String(0), args.Error(1)
}

func (walletService *walletServiceMock) Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error) {
	args := walletService.Called(recipientWalletId, amount)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (walletService *walletServiceMock) Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error) {
	args := walletService.Called(senderWalletId, recipientWalletId, amount)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (walletService *walletServiceMock) GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error) {
	args := walletService.Called(limit, offset, filter)
	return args.Get(0).([]*model.Transaction), args.Error(1)
}

func TestCreateWallet(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	reqBody := dto.CreateWalletRequest{
		Name:     "wallet",
		Currency: "USD",
	}
	reqBodyBuf := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBuf).Encode(reqBody); err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest("POST", "/wallets", reqBodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	walletService.On("CreateWallet", model.Wallet{
		Name:     "wallet",
		Currency: "USD",
	}).Return("1001", nil)

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	respBody := new(dto.CreateWalletResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1001", respBody.ID)

	walletService.AssertNumberOfCalls(t, "CreateWallet", 1)
	walletService.AssertExpectations(t)
}

func TestCreateWalletCurrencyValidation(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	reqBody := dto.CreateWalletRequest{
		Name:     "wallet",
		Currency: "EUR",
	}
	reqBodyBuf := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBuf).Encode(reqBody); err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest("POST", "/wallets", reqBodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	respBody := new(dto.ErrorResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "invalid request body", respBody.Message)

	walletService.AssertNumberOfCalls(t, "CreateWallet", 0)
	walletService.AssertExpectations(t)
}

func TestDeposit(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	reqBody := dto.DepositRequest{
		Amount: 10000,
	}
	reqBodyBuf := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBuf).Encode(reqBody); err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest("POST", "/wallets/1001/deposit", reqBodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	walletService.On(
		"Deposit", "1001", uint64(10000),
	).Return(
		&model.Transaction{
			Amount:      10000,
			ProcessedAt: time.Now().UTC(),
			RecipientWallet: model.Wallet{
				ID:      "1001",
				Balance: 10000,
			},
			OperationType: model.Deposit,
		}, nil,
	)

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	respBody := new(dto.DepositResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(10000), respBody.Balance)

	walletService.AssertNumberOfCalls(t, "Deposit", 1)
	walletService.AssertExpectations(t)
}

func TestTransfer(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	reqBody := dto.TransferRequest{
		Amount:         10000,
		SenderWalletId: "1002",
	}
	reqBodyBuf := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBuf).Encode(reqBody); err != nil {
		t.Fatal(err)
	}
	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest("POST", "/wallets/1001/transfer", reqBodyBuf)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	walletService.On(
		"Transfer", "1002", "1001", uint64(10000),
	).Return(
		&model.Transaction{
			Amount:      10000,
			ProcessedAt: time.Now().UTC(),
			SenderWallet: &model.Wallet{
				ID:      "1002",
				Balance: 0,
			},
			RecipientWallet: model.Wallet{
				ID:      "1001",
				Balance: 10000,
			},
			OperationType: model.Transfer,
		}, nil,
	)

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	respBody := new(dto.TransferResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(0), respBody.SenderWalletBalance)
	assert.Equal(t, uint64(10000), respBody.Balance)

	walletService.AssertNumberOfCalls(t, "Transfer", 1)
	walletService.AssertExpectations(t)
}

func TestGetDepositTransactionsJson(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest(
		"GET",
		"/wallets/1001/transactions?limit=10&offset=0&operation_type=deposit&processed_at.gte=2022-01-10T00:00:00Z&processed_at.lte=2022-01-10T23:59:59Z",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	processedAtGte, err := time.Parse(time.RFC3339Nano, "2022-01-10T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	processedAtLte, err := time.Parse(time.RFC3339Nano, "2022-01-10T23:59:59Z")
	if err != nil {
		t.Fatal(err)
	}
	walletService.On(
		"GetTransactions", 10, 0, model.TransactionFilter{
			WalletId:       "1001",
			ProcessedAtGte: processedAtGte,
			ProcessedAtLte: processedAtLte,
			OperationType:  model.Deposit,
		},
	).Return(
		[]*model.Transaction{
			{
				Amount:      10000,
				ProcessedAt: time.Now().UTC(),
				RecipientWallet: model.Wallet{
					ID:      "1001",
					Balance: 20000,
				},
				OperationType: model.Deposit,
			},
		}, nil,
	)

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	respBody := new(dto.TransactionsResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	transaction := (*respBody)[0]
	assert.Equal(t, "deposit", transaction.OperationType)
	assert.Equal(t, uint64(10000), transaction.Amount)
	assert.Nil(t, transaction.SenderWalletId)
	assert.Nil(t, transaction.SenderWalletBalance)
	assert.False(t, transaction.SenderWalletMe)
	assert.Nil(t, transaction.RecipientWalletId)
	assert.Nil(t, transaction.RecipientWalletBalance)
	assert.False(t, transaction.RecipientWalletMe)
	assert.Equal(t, uint64(20000), transaction.Balance)
	assert.NotNil(t, transaction.ProcessedAt)

	walletService.AssertNumberOfCalls(t, "GetTransactions", 1)
	walletService.AssertExpectations(t)
}

func TestGetTransferTransactionsJson(t *testing.T) {
	// given
	walletService := new(walletServiceMock)

	router := mux.NewRouter()
	NewWalletsApi(logger, router, walletService)
	req, err := http.NewRequest(
		"GET",
		"/wallets/1001/transactions?limit=10&offset=0&operation_type=transfer&processed_at.gte=2022-01-10T00:00:00Z&processed_at.lte=2022-01-10T23:59:59Z",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	processedAtGte, err := time.Parse(time.RFC3339Nano, "2022-01-10T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	processedAtLte, err := time.Parse(time.RFC3339Nano, "2022-01-10T23:59:59Z")
	if err != nil {
		t.Fatal(err)
	}
	walletService.On(
		"GetTransactions", 10, 0, model.TransactionFilter{
			WalletId:       "1001",
			ProcessedAtGte: processedAtGte,
			ProcessedAtLte: processedAtLte,
			OperationType:  model.Transfer,
		},
	).Return(
		[]*model.Transaction{
			{
				Amount:      10000,
				ProcessedAt: time.Now().UTC(),
				SenderWallet: &model.Wallet{
					ID:      "1002",
					Balance: 20000,
				},
				RecipientWallet: model.Wallet{
					ID:      "1001",
					Balance: 30000,
				},
				OperationType: model.Transfer,
			},
		}, nil,
	)

	// when
	router.ServeHTTP(recorder, req)

	// then
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	respBody := new(dto.TransactionsResponse)
	if err := json.Unmarshal(recorder.Body.Bytes(), respBody); err != nil {
		t.Fatal(err)
	}
	transaction := (*respBody)[0]
	assert.Equal(t, "transfer", transaction.OperationType)
	assert.Equal(t, uint64(10000), transaction.Amount)
	assert.Equal(t, "1002", *transaction.SenderWalletId)
	assert.Equal(t, uint64(20000), *transaction.SenderWalletBalance)
	assert.False(t, transaction.SenderWalletMe)
	assert.Nil(t, transaction.RecipientWalletId)
	assert.Nil(t, transaction.RecipientWalletBalance)
	assert.True(t, transaction.RecipientWalletMe)
	assert.Equal(t, uint64(30000), transaction.Balance)
	assert.NotNil(t, transaction.ProcessedAt)

	walletService.AssertNumberOfCalls(t, "GetTransactions", 1)
	walletService.AssertExpectations(t)
}
