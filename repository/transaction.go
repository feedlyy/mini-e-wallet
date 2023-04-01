package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) GetsByUsrID(ctx context.Context, usrID string) ([]domain.Transaction, error) {
	var (
		res  []domain.Transaction
		err  error
		stmt *sqlx.Stmt
		rows *sqlx.Rows
		sql  string
	)
	sql, _, err = sq.Select("id", "status", "transaction_at", "type", "amount", "reference_id").
		From("transactions").Where(sq.And{
		sq.Eq{"transaction_by": "id"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Transactions - Repository|err when generate sql, err:%v", err)
		return []domain.Transaction{}, err
	}

	stmt, err = t.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Transactions - Repository|err when init prepare statement, err:%v", err)
		return []domain.Transaction{}, err
	}
	defer stmt.Close()

	rows, err = stmt.QueryxContext(ctx, usrID)
	if err != nil {
		logrus.Errorf("Transactions - Repository|err when get data, err:%v", err)
		return []domain.Transaction{}, err
	}

	for rows.Next() {
		var tmp domain.Transaction
		if err = rows.Scan(&tmp.Id, &tmp.Status, &tmp.TransactionAt, &tmp.Type, &tmp.Amount, &tmp.ReferenceID); err != nil {
			logrus.Errorf("Transactions - Repository|err when scan data, err:%v", err)
			return []domain.Transaction{}, err
		}
		res = append(res, tmp)
	}

	return res, nil
}

func (t *transactionRepository) Store(ctx context.Context, transaction domain.Transaction) error {
	var (
		err error
		sql string
	)
	sql, _, err = sq.Insert("transactions").Columns("id", "status", "type", "amount", "reference_id", "transaction_at", "transaction_by").
		Values("id", "status", "type", "amount", "reference_id", "transaction_at", "transaction_by").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Transactions - Repository|err when generate sql, err:%v", err)
		return err
	}

	_, err = t.db.ExecContext(ctx, sql, transaction.Id, transaction.Status, transaction.Type, transaction.Amount, transaction.ReferenceID,
		transaction.TransactionAt, transaction.TransactionBy)
	if err != nil {
		logrus.Errorf("Transactions - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}
