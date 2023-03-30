package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"time"
)

type accountService struct {
	accountRepo domain.AccountRepository
	tokenRepo   domain.TokenRepository
}

func NewAccountService(a domain.AccountRepository, t domain.TokenRepository) domain.AccountService {
	return &accountService{
		accountRepo: a,
		tokenRepo:   t,
	}
}

func (a *accountService) Register(ctx context.Context, customerXID string) (string, error) {
	var (
		err          error
		token, usrID string
		tx           *sqlx.Tx
		tokenData    domain.Tokens
		loc          *time.Location
		expiredAt    time.Time
	)
	token = helpers.GenerateRandomUUID()

	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Errorf("Account - Service|Err when get location %v", err)
		return "", err
	}
	expiredAt = time.Now().In(loc)

	// begin tx
	tx, err = a.accountRepo.BeginTx(ctx)
	if err != nil {
		logrus.Errorf("Account - Service|err when initiate tx, err:%v", err)
		return "", err
	}

	// Rollback the transaction if there's an error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// create account
	usrID, err = a.accountRepo.Store(ctx, customerXID, tx)
	if err != nil {
		return "", err
	}

	// create token
	tokenData = domain.Tokens{
		AccountID:  usrID,
		Token:      token,
		Expiration: expiredAt.Add(time.Hour),
		CreatedAt:  time.Now().In(loc),
	}
	err = a.tokenRepo.Store(ctx, tokenData, tx)
	if err != nil {
		return "", err
	}

	// Commit the transaction if everything is successful
	if err = tx.Commit(); err != nil {
		logrus.Errorf("Account - Service|err when commit tx, err:%v", err)
		return "", err
	}

	return token, nil
}
