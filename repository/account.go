package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
)

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) domain.AccountRepository {
	return &accountRepository{db: db}
}

func (a *accountRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return a.db.Beginx()
}

func (a *accountRepository) Store(ctx context.Context, customerXID string, tx *sqlx.Tx) (string, error) {
	var (
		err error
		sql string
		id  = helpers.GenerateRandomUUID()
	)
	sql, _, err = sq.Insert("accounts").Columns("id", "customer_xid").
		Values("id", "customer").PlaceholderFormat(sq.Dollar).ToSql()

	if tx == nil {
		_, err = a.db.ExecContext(ctx, sql, id, customerXID)
	} else {
		_, err = tx.ExecContext(ctx, sql, id, customerXID)
	}

	if err != nil {
		logrus.Errorf("Accounts - Repository|err when store data, err:%v", err)
		return "", err
	}

	return id, nil
}
