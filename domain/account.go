package domain

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Accounts struct {
	Id          string `json:"id"`
	CustomerXID string `json:"customer_xid"`
}

type AccResp struct {
	Token string `json:"token,omitempty"`
}

type AccountRepository interface {
	Store(ctx context.Context, customerXID string, tx *sqlx.Tx) (string, error)
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type AccountService interface {
	Register(ctx context.Context, customerXID string) (string, error)
}
