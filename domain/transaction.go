package domain

import (
	"context"
	"time"
)

type Transaction struct {
	Id            string    `json:"id,omitempty"`
	Status        string    `json:"status"`
	TransactionAt time.Time `json:"transacted_at"`
	Type          string    `json:"type"`
	Amount        int       `json:"amount"`
	ReferenceID   string    `json:"reference_id"`
	TransactionBy string    `json:"-"`
}

type Deposit struct {
	Id          string    `json:"id,omitempty"`
	DepositedBy string    `json:"deposited_by,omitempty"`
	Status      string    `json:"status,omitempty"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      int       `json:"amount,omitempty"`
	ReferenceId string    `json:"reference_id,omitempty"`
}

type TransactionRepository interface {
	GetsByUsrID(ctx context.Context, usrID string) ([]Transaction, error)
	Store(ctx context.Context, transaction Transaction) error
}
