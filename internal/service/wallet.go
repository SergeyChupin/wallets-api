package service

import (
	"errors"

	"github.com/SergeyChupin/wallets-api/internal/model"
	"github.com/SergeyChupin/wallets-api/internal/repository"
)

type WalletService interface {
	CreateWallet(wallet model.Wallet) (string, error)
	Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error)
	Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error)
	GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error)
}

type walletService struct {
	walletRepository repository.WalletRepository
}

func NewWalletService(walletRepository repository.WalletRepository) *walletService {
	return &walletService{
		walletRepository: walletRepository,
	}
}

func (walletService *walletService) CreateWallet(wallet model.Wallet) (string, error) {
	return walletService.walletRepository.CreateWallet(wallet)
}

func (walletService *walletService) Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error) {
	return walletService.walletRepository.Deposit(recipientWalletId, amount)
}

func (walletService *walletService) Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error) {
	if senderWalletId == recipientWalletId {
		return nil, errors.New("should be specified different wallets")
	}
	return walletService.walletRepository.Transfer(senderWalletId, recipientWalletId, amount)
}

func (walletService *walletService) GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error) {
	return walletService.walletRepository.GetTransactions(limit, offset, filter)
}
