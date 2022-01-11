package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/SergeyChupin/wallets-api/internal/model"
)

type WalletRepository interface {
	CreateWallet(wallet model.Wallet) (string, error)
	Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error)
	Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error)
	GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *walletRepository {
	return &walletRepository{
		db: db,
	}
}

func (walletRepository *walletRepository) CreateWallet(wallet model.Wallet) (string, error) {
	var id string
	if err := walletRepository.db.QueryRow(
		"INSERT INTO wallets(name, currency, balance) VALUES($1, $2, $3) RETURNING id",
		wallet.Name,
		wallet.Currency,
		0,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("WalletRepository - CreateWallet - walletRepository.db.QueryRow: %w", err)
	}
	return id, nil
}

func (walletRepository *walletRepository) Deposit(recipientWalletId string, amount uint64) (*model.Transaction, error) {
	tx, err := walletRepository.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("WalletRepository - Deposit - walletRepository.db.Begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()

	var recipientWalletBalance uint64
	if err = tx.QueryRow(
		"UPDATE wallets SET balance = balance + $1 WHERE id = $2 RETURNING balance",
		amount,
		recipientWalletId,
	).Scan(&recipientWalletBalance); err != nil {
		return nil, fmt.Errorf("WalletRepository - Deposit - tx.QueryRow: %w", err)
	}

	if _, err = tx.Exec(
		"INSERT INTO transactions(operation_type, amount, recipient_wallet_id, recipient_wallet_balance, processed_at) VALUES($1, $2, $3, $4, $5)",
		model.Deposit,
		amount,
		recipientWalletId,
		recipientWalletBalance,
		now,
	); err != nil {
		return nil, fmt.Errorf("WalletRepository - Deposit - tx.Exec: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("WalletRepository - Deposit - tx.Commit: %w", err)
	}

	return &model.Transaction{
		Amount:      amount,
		ProcessedAt: now,
		RecipientWallet: model.Wallet{
			ID:      recipientWalletId,
			Balance: recipientWalletBalance,
		},
		OperationType: model.Deposit,
	}, nil
}

func (walletRepository *walletRepository) Transfer(senderWalletId string, recipientWalletId string, amount uint64) (*model.Transaction, error) {
	tx, err := walletRepository.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("WalletRepository - Transfer - walletRepository.db.Begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()

	var senderWalletBalance uint64
	if err = tx.QueryRow(
		"UPDATE wallets SET balance = balance - $1 WHERE id = $2 RETURNING balance",
		amount,
		senderWalletId,
	).Scan(&senderWalletBalance); err != nil {
		return nil, fmt.Errorf("WalletRepository - Transfer - tx.QueryRow: %w", err)
	}

	var recipientWalletBalance uint64
	if err = tx.QueryRow(
		"UPDATE wallets SET balance = balance + $1 WHERE id = $2 RETURNING balance",
		amount,
		recipientWalletId,
	).Scan(&recipientWalletBalance); err != nil {
		return nil, fmt.Errorf("WalletRepository - Transfer - tx.QueryRow: %w", err)
	}

	if _, err = tx.Exec(
		"INSERT INTO transactions(operation_type, amount, sender_wallet_id, sender_wallet_balance, recipient_wallet_id, recipient_wallet_balance, processed_at) VALUES($1, $2, $3, $4, $5, $6, $7)",
		model.Transfer,
		amount,
		senderWalletId,
		senderWalletBalance,
		recipientWalletId,
		recipientWalletBalance,
		now,
	); err != nil {
		return nil, fmt.Errorf("WalletRepository - Transfer - tx.Exec: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("WalletRepository - Transfer - tx.Commit: %w", err)
	}

	return &model.Transaction{
		Amount:      amount,
		ProcessedAt: now,
		SenderWallet: &model.Wallet{
			ID:      senderWalletId,
			Balance: senderWalletBalance,
		},
		RecipientWallet: model.Wallet{
			ID:      recipientWalletId,
			Balance: recipientWalletBalance,
		},
		OperationType: model.Transfer,
	}, nil
}

type transaction struct {
	amount                 uint64
	processedAt            time.Time
	senderWalletId         sql.NullString
	senderWalletBalance    sql.NullString
	recipientWalletId      string
	recipientWalletBalance uint64
	operationType          string
}

func (walletRepository *walletRepository) GetTransactions(limit int, offset int, filter model.TransactionFilter) ([]*model.Transaction, error) {
	query := "SELECT operation_type, amount, sender_wallet_id, sender_wallet_balance, recipient_wallet_id, recipient_wallet_balance, processed_at FROM transactions WHERE (sender_wallet_id = $1 OR recipient_wallet_id = $1)"
	var filterValues []interface{}
	filterValues = append(filterValues, filter.WalletId)
	if filter.OperationType != model.UnknownOperation {
		filterValues = append(filterValues, filter.OperationType)
		query += " AND operation_type = $" + strconv.Itoa(len(filterValues))
	}
	if !filter.ProcessedAtGte.IsZero() {
		filterValues = append(filterValues, filter.ProcessedAtGte)
		query += " AND processed_at >= $" + strconv.Itoa(len(filterValues))
	}
	if !filter.ProcessedAtLte.IsZero() {
		filterValues = append(filterValues, filter.ProcessedAtLte)
		query += " AND processed_at <= $" + strconv.Itoa(len(filterValues))
	}
	query += " ORDER BY processed_at DESC"
	if limit > -1 {
		filterValues = append(filterValues, limit)
		query += " LIMIT $" + strconv.Itoa(len(filterValues))
	}
	if offset > -1 {
		filterValues = append(filterValues, offset)
		query += " OFFSET $" + strconv.Itoa(len(filterValues))
	}

	rows, err := walletRepository.db.Query(
		query, filterValues...,
	)
	if err != nil {
		return nil, fmt.Errorf("WalletRepository - GetTransactions - walletRepository.db.Query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var transactions []*model.Transaction
	for rows.Next() {
		transactionEntity := transaction{}
		if err = rows.Scan(
			&transactionEntity.operationType,
			&transactionEntity.amount,
			&transactionEntity.senderWalletId,
			&transactionEntity.senderWalletBalance,
			&transactionEntity.recipientWalletId,
			&transactionEntity.recipientWalletBalance,
			&transactionEntity.processedAt,
		); err != nil {
			return nil, fmt.Errorf("WalletRepository - GetTransactions - rows.Scan: %w", err)
		}
		operationType, err := model.FromString(transactionEntity.operationType)
		if err != nil {
			return nil, fmt.Errorf("WalletRepository - GetTransactions - model.FromString: %w", err)
		}
		transaction := new(model.Transaction)
		transaction.Amount = transactionEntity.amount
		transaction.ProcessedAt = transactionEntity.processedAt
		transaction.RecipientWallet.ID = transactionEntity.recipientWalletId
		transaction.RecipientWallet.Balance = transactionEntity.recipientWalletBalance
		transaction.OperationType = operationType
		if transactionEntity.senderWalletBalance.Valid {
			senderWalletBalance, err := strconv.ParseUint(transactionEntity.senderWalletBalance.String, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("WalletRepository - GetTransactions - strconv.ParseUint: %w", err)
			}
			transaction.SenderWallet = &model.Wallet{
				ID:      transactionEntity.senderWalletId.String,
				Balance: senderWalletBalance,
			}
		}
		transactions = append(transactions, transaction)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("WalletRepository - GetTransactions - rows.Err: %w", err)
	}

	return transactions, nil
}
