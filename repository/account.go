package repository

import (
	"context"
	sql2 "database/sql"
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
	if err != nil {
		logrus.Errorf("Accounts - Repository|err when generate sql, err:%v", err)
		return "", err
	}

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

func (a *accountRepository) GetByCustID(ctx context.Context, customerXID string) (domain.Accounts, error) {
	var (
		err  error
		res  = domain.Accounts{}
		sql  string
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("id", "customer_xid").From("accounts").Where(sq.And{
		sq.Eq{"customer_xid": "custID"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Accounts - Repository|err when generate sql, err:%v", err)
		return domain.Accounts{}, err
	}

	stmt, err = a.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Accounts - Repository|err when init prepare statement, err:%v", err)
		return domain.Accounts{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, customerXID)
	err = row.Scan(&res.Id, &res.CustomerXID)
	if err != nil && err != sql2.ErrNoRows {
		logrus.Errorf("Accounts - Repository|err when scan, err:%v", err)
		return domain.Accounts{}, err
	}

	return res, nil
}
