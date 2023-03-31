package domain

import (
	"context"
	"database/sql"
	"time"
)

type Wallets struct {
	Id         string       `json:"id,omitempty"`
	OwnedBy    string       `json:"owned_by,omitempty"`
	Status     string       `json:"status"`
	EnableAt   time.Time    `json:"enable_at"`
	Balance    int          `json:"balance"`
	DisabledAt sql.NullTime `json:"-"`
}

type WalletRepository interface {
	Store(ctx context.Context, usrID string) error
	GetByOwnedID(ctx context.Context, id string) (Wallets, error)
	Update(ctx context.Context, wallet Wallets) error
}

type WalletService interface {
	Enable(ctx context.Context, token string) (Wallets, error)
	Disable(ctx context.Context, token string) (Wallets, error)
	Balance(ctx context.Context, token string) (Wallets, error)
}
