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
	walletRepo domain.WalletRepository
	tokenRepo  domain.TokenRepository
}

func NewWalletService(w domain.WalletRepository, t domain.TokenRepository) domain.WalletService {
	return &walletService{
		walletRepo: w,
		tokenRepo:  t,
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

	if walletAcc == (domain.Wallets{}) || walletAcc.Status == helpers.DisabledStatus {
		err = errors.New(helpers.ErrWalletNotExists)
		logrus.Errorf("Wallet - Service|Err data alr enabled %v", err)
		return walletAcc, err
	}

	return walletAcc, nil
}
