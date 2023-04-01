package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"net/http"
	"strings"
)

type Middleware struct {
	tokenRepo  domain.TokenRepository
	walletRepo domain.WalletRepository
}

func NewMiddleware(t domain.TokenRepository, w domain.WalletRepository) Middleware {
	handler := &Middleware{
		tokenRepo:  t,
		walletRepo: w,
	}

	return *handler
}

// AuthMiddleware function to authenticate requests using the provided token
func (m *Middleware) AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Error("Middleware | Empty auth header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "token" {
			logrus.Error("Middleware | Empty token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := splitHeader[1]

		checkToken, err := m.tokenRepo.GetByToken(ctx, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if checkToken == (domain.Tokens{}) {
			logrus.Error("Middleware | Empty data")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the token to the request context
		ctx = context.WithValue(ctx, "token", token)

		// If authentication succeeded, call the next handler with the modified context
		next(w, r.WithContext(ctx), ps)
	}
}

// WalletMiddleware function to authenticate requests using the provided token and checking its wallet
func (m *Middleware) WalletMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.Background()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Error("Middleware | Empty auth header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		// Extract the token from the header
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "token" {
			logrus.Error("Middleware | Empty token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := splitHeader[1]

		checkToken, err := m.tokenRepo.GetByToken(ctx, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if checkToken == (domain.Tokens{}) {
			logrus.Error("Middleware | Empty data")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// get data wallet
		walletAcc, err := m.walletRepo.GetByOwnedID(ctx, checkToken.AccountID)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if walletAcc == (domain.Wallets{}) || walletAcc.Status == helpers.DisabledStatus {
			resp := helpers.Response{
				Status: helpers.FailMsg,
				Data:   helpers.ErrResp{Err: helpers.ErrWalletDisabled},
			}
			logrus.Error("Middleware | Err data alr disabled")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Add the token to the request context
		ctx = context.WithValue(ctx, "token", token)

		// If authentication succeeded, call the next handler with the modified context
		next(w, r.WithContext(ctx), ps)
	}
}
