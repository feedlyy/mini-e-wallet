package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"time"
)

type walletService struct {
	walletRepo      domain.WalletRepository
	tokenRepo       domain.TokenRepository
	transactionRepo domain.TransactionRepository
	kafka           domain.KafkaProducer
}

func NewWalletService(w domain.WalletRepository, t domain.TokenRepository, tr domain.TransactionRepository, k domain.KafkaProducer) domain.WalletService {
	return &walletService{
		walletRepo:      w,
		tokenRepo:       t,
		transactionRepo: tr,
		kafka:           k,
	}
}

func (w *walletService) Enable(ctx context.Context, token string) (domain.Wallets, error) {
	var (
		err       error
		tokenData domain.Tokens
		walletAcc domain.Wallets
	)
	// get id by tokens
	tokenData, err = w.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.Wallets{}, err
	}

	// get data wallet if exists
	walletAcc, err = w.walletRepo.GetByOwnedID(ctx, tokenData.AccountID)
	if err != nil && err != sql.ErrNoRows {
		return domain.Wallets{}, err
	}

	if walletAcc != (domain.Wallets{}) {
		if walletAcc.Status == helpers.DisabledStatus {
			// update again status to enabled
			walletAcc.Status = helpers.EnabledStatus
			walletAcc.DisabledAt = sql.NullTime{Time: time.Time{}}
			err = w.walletRepo.Update(ctx, walletAcc)
			if err != nil {
				return domain.Wallets{}, err
			}

			return walletAcc, nil
		}

		if walletAcc.Status == helpers.EnabledStatus {
			err = errors.New(helpers.ErrAlreadyEnabled)
			logrus.Errorf("Wallet - Service|Err data alr enabled %v", err)
			return walletAcc, err
		}
	}

	// store new wallet
	err = w.walletRepo.Store(ctx, tokenData.AccountID)
	if err != nil {
		return domain.Wallets{}, err
	}

	// get for returned inserted id
	walletAcc, err = w.walletRepo.GetByOwnedID(ctx, tokenData.AccountID)
	if err != nil {
		return domain.Wallets{}, err
	}

	return walletAcc, nil
}

func (w *walletService) Disable(ctx context.Context, token string, isDisabled bool) (domain.Wallets, error) {
	var (
		err       error
		tokenData domain.Tokens
		walletAcc domain.Wallets
	)
	// get id by tokens
	tokenData, err = w.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.Wallets{}, err
	}

	// get data wallet
	walletAcc, err = w.walletRepo.GetByOwnedID(ctx, tokenData.AccountID)
	if err != nil && err != sql.ErrNoRows {
		return domain.Wallets{}, err
	}

	// update status to disabled
	if isDisabled {
		walletAcc.Status = helpers.DisabledStatus
	}
	walletAcc.DisabledAt = sql.NullTime{Time: time.Now(), Valid: true}
	err = w.walletRepo.Update(ctx, walletAcc)
	if err != nil {
		return domain.Wallets{}, err
	}

	return walletAcc, nil
}

func (w *walletService) Balance(ctx context.Context, token string) (domain.Wallets, error) {
	var (
		err       error
		tokenData domain.Tokens
		walletAcc domain.Wallets
	)
	// get id by tokens
	tokenData, err = w.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.Wallets{}, err
	}

	// get data wallet
	walletAcc, err = w.walletRepo.GetByOwnedID(ctx, tokenData.AccountID)
	if err != nil && err != sql.ErrNoRows {
		return domain.Wallets{}, err
	}

	return walletAcc, nil
}

func (w *walletService) Transactions(ctx context.Context, token string) ([]domain.Transaction, error) {
	var (
		err       error
		tokenData domain.Tokens
		res       []domain.Transaction
	)
	// get id by tokens
	tokenData, err = w.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return []domain.Transaction{}, err
	}

	res, err = w.transactionRepo.GetsByUsrID(ctx, tokenData.AccountID)
	if err != nil {
		return []domain.Transaction{}, err
	}

	return res, nil
}

func (w *walletService) AddFunds(ctx context.Context, token string, transaction domain.Transaction) (domain.Transaction, error) {
	var (
		err       error
		tokenData domain.Tokens
		id        = helpers.GenerateRandomUUID()
	)

	// get id by tokens
	tokenData, err = w.tokenRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.Transaction{}, err
	}

	transaction.Id = id
	transaction.TransactionAt = time.Now()
	transaction.Status = helpers.SuccessMsg
	transaction.TransactionBy = tokenData.AccountID
	transaction.Type = helpers.Deposit

	// store transaction
	err = w.transactionRepo.Store(ctx, transaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	// update wallet
	err = w.kafka.SendMessage(helpers.WalletTopic, transaction)
	if err != nil {
		logrus.Errorf("Wallet - Service|err send email user with kafka, err:%v", err)
		return domain.Transaction{}, err
	}

	return transaction, nil
}
