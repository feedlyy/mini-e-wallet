package handler

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"net/http"
	"time"
)

type WalletHandler struct {
	walletService domain.WalletService
	timeout       time.Duration
}

func NewWalletHandler(w domain.WalletService, timeout time.Duration) WalletHandler {
	handler := &WalletHandler{
		walletService: w,
		timeout:       timeout,
	}

	return *handler
}

func (o *WalletHandler) EnableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err   error
		token = r.Context().Value("token").(string)
		resp  = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		res domain.Wallets
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	res, err = o.walletService.Enable(ctx, token)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = res
	json.NewEncoder(w).Encode(resp)
	return
}

func (o *WalletHandler) DisableWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err   error
		token = r.Context().Value("token").(string)
		resp  = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		res domain.Wallets
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	res, err = o.walletService.Disable(ctx, token)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = res
	json.NewEncoder(w).Encode(resp)
	return
}
