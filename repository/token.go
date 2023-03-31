package repository

import (
	"context"
	sql2 "database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"time"
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
	if err != nil {
		logrus.Errorf("Tokens - Repository|err when generate sql, err:%v", err)
		return err
	}

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

func (a *tokenRepository) GetByToken(ctx context.Context, token string) (domain.Tokens, error) {
	var (
		err  error
		res  = domain.Tokens{}
		sql  string
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("account_id").From("tokens").Where(sq.And{
		sq.Eq{"token": token},
		sq.GtOrEq{"expiration": time.Now()},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Tokens - Repository|err when generate sql, err:%v", err)
		return domain.Tokens{}, err
	}

	stmt, err = a.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Tokens - Repository|err when init prepare statement, err:%v", err)
		return domain.Tokens{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, token, time.Now())
	err = row.Scan(&res.AccountID)
	if err != nil && err != sql2.ErrNoRows {
		logrus.Errorf("Tokens - Repository|err when get by token, err:%v", err)
		return domain.Tokens{}, err
	}

	return res, nil
}
