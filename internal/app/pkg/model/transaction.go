package model

import (
	"errors"
	"time"
)

type OperationType struct {
	value string
}

func (operationType OperationType) String() string {
	return operationType.value
}

var (
	UnknownOperation = OperationType{""}
	Deposit          = OperationType{"deposit"}
	Transfer         = OperationType{"transfer"}
)

func FromString(value string) (OperationType, error) {
	switch value {
	case Deposit.value:
		return Deposit, nil
	case Transfer.value:
		return Transfer, nil
	}
	return UnknownOperation, errors.New("unknown operation type")
}

type Transaction struct {
	Amount          uint64
	ProcessedAt     time.Time
	SenderWallet    *Wallet
	RecipientWallet Wallet
	OperationType   OperationType
}

type TransactionFilter struct {
	WalletId       string
	ProcessedAtGte time.Time
	ProcessedAtLte time.Time
	OperationType  OperationType
}
