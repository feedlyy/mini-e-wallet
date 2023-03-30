package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
)

type tokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) domain.TokenRepository {
	return &tokenRepository{db: db}
}

func (a *tokenRepository) Store(ctx context.Context, token domain.Tokens, tx *sqlx.Tx) error {
	var (
		err error
		sql string
	)
	sql, _, err = sq.Insert("tokens").Columns("account_id", "token", "expiration", "created_at").
		Values("account_id", "token", "expiration", "created_at").
		PlaceholderFormat(sq.Dollar).ToSql()

	if tx == nil {
		_, err = a.db.ExecContext(ctx, sql, token.AccountID, token.Token, token.Expiration, token.CreatedAt)
	} else {
		_, err = tx.ExecContext(ctx, sql, token.AccountID, token.Token, token.Expiration, token.CreatedAt)
	}
	if err != nil {
		logrus.Errorf("Tokens - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}
