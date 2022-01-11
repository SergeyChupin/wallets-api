package service

import (
	"errors"
	"fmt"

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
	id, err := walletService.walletRepository.CreateWallet(wallet)
	if err != nil {
		return "", fmt.Errorf("WalletService - CreateWallet - walletService.walletRepository.CreateWallet: %w", err)
	}
	return id, nil
}

func (walletService *walletService) Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error) {
	transaction, err := walletService.walletRepository.Deposit(recipientWalletId, amount)
	if err != nil {
		return nil, fmt.Errorf("WalletService - Deposit - walletService.walletRepository.Deposit: %w", err)
	}
	return transaction, nil
}

func (walletService *walletService) Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error) {
	if senderWalletId == recipientWalletId {
		return nil, errors.New("WalletService - Transfer - sender wallet should be different than recipient wallet")
	}
	transaction, err := walletService.walletRepository.Transfer(senderWalletId, recipientWalletId, amount)
	if err != nil {
		return nil, fmt.Errorf("WalletService - Transfer - walletService.walletRepository.Transfer: %w", err)
	}
	return transaction, nil
}

func (walletService *walletService) GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error) {
	transactions, err := walletService.walletRepository.GetTransactions(limit, offset, filter)
	if err != nil {
		return nil, fmt.Errorf("WalletService - GetTransactions - walletService.walletRepository.GetTransactions: %w", err)
	}
	return transactions, nil
}
