package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"time"
)

type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (a *walletRepository) Store(ctx context.Context, usrID string) error {
	var (
		err error
		sql string
		id  = helpers.GenerateRandomUUID()
	)
	sql, _, err = sq.Insert("wallets").Columns("id", "owned_by", "status", "enabled_at", "balance").
		Values("id", "owned_by", "status", "enabled_at", "balance").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when generate sql, err:%v", err)
		return err
	}

	_, err = a.db.ExecContext(ctx, sql, id, usrID, helpers.EnabledStatus, time.Now(), 0)
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}

func (a *walletRepository) GetByOwnedID(ctx context.Context, id string) (domain.Wallets, error) {
	var (
		err  error
		sql  string
		res  domain.Wallets
		stmt *sqlx.Stmt
		row  *sqlx.Row
	)
	sql, _, err = sq.Select("id", "owned_by", "status", "enabled_at", "balance", "disabled_at").From("wallets").Where(sq.And{
		sq.Eq{"owned_by": "id"},
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when generate sql, err:%v", err)
		return domain.Wallets{}, err
	}

	stmt, err = a.db.PreparexContext(ctx, sql)
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when init prepare statement, err:%v", err)
		return domain.Wallets{}, err
	}
	defer stmt.Close()

	row = stmt.QueryRowxContext(ctx, id)
	err = row.Scan(&res.Id, &res.OwnedBy, &res.Status, &res.EnableAt, &res.Balance, &res.DisabledAt)
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when get data, err:%v", err)
		return domain.Wallets{}, err
	}

	return res, nil
}

func (a *walletRepository) Update(ctx context.Context, wallet domain.Wallets) error {
	var (
		err error
		sql string
	)
	sql, _, err = sq.Update("wallets").
		Set("disabled_at", wallet.DisabledAt).
		Set("status", wallet.Status).
		Set("balance", wallet.Balance).
		Where(sq.Eq{"id": "id"}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when generate sql, err:%v", err)
		return err
	}

	_, err = a.db.ExecContext(ctx, sql, wallet.DisabledAt, wallet.Status, wallet.Balance, wallet.Id)
	if err != nil {
		logrus.Errorf("Wallets - Repository|err when update data, err:%v", err)
		return err
	}

	return nil
}
