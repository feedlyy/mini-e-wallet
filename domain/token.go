package domain

import (
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type Tokens struct {
	Id         string    `json:"id,omitempty"`
	AccountID  string    `json:"account_id,omitempty"`
	Token      string    `json:"token,omitempty"`
	Expiration time.Time `json:"-"`
	CreatedAt  time.Time `json:"-"`
}

type TokenRepository interface {
	Store(ctx context.Context, token Tokens, tx *sqlx.Tx) error
	GetByToken(ctx context.Context, token string) (Tokens, error)
	Update(ctx context.Context, tokens Tokens, tx *sqlx.Tx) error
}
