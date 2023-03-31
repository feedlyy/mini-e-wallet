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
}

type TransactionRepository interface {
	GetsByUsrID(ctx context.Context, usrID string) ([]Transaction, error)
}
