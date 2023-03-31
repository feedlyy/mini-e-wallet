package handler

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mini-e-wallet/domain"
	"mini-e-wallet/helpers"
	"net/http"
	"strconv"
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
		res     domain.Wallets
		errResp = helpers.ErrResp{}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	res, err = o.walletService.Enable(ctx, token)
	if err != nil {
		errResp.Err = err.Error()
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		switch {
		case err.Error() == helpers.ErrAlreadyEnabled:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
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
		res            domain.Wallets
		formIsDisabled = r.PostFormValue("is_disabled")
		isDisabled     bool
		errResp        = helpers.ErrResp{}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	if formIsDisabled == "" {
		errResp.Err = "Missing required field: is_disabled"
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	isDisabled, err = strconv.ParseBool(formIsDisabled)
	if err != nil {
		errResp.Err = err.Error()
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res, err = o.walletService.Disable(ctx, token, isDisabled)
	if err != nil {
		errResp.Err = err.Error()
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = res
	json.NewEncoder(w).Encode(resp)
	return
}

func (o *WalletHandler) ViewBalance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err   error
		token = r.Context().Value("token").(string)
		resp  = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		res     domain.Wallets
		errResp = helpers.ErrResp{}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	res, err = o.walletService.Balance(ctx, token)
	if err != nil {
		errResp.Err = err.Error()
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		switch {
		case err.Error() == helpers.ErrWalletDisabled:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = res
	json.NewEncoder(w).Encode(resp)
	return
}

func (o *WalletHandler) ListTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err   error
		token = r.Context().Value("token").(string)
		resp  = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		res     []domain.Transaction
		errResp = helpers.ErrResp{}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	res, err = o.walletService.Transactions(ctx, token)
	if err != nil {
		errResp.Err = err.Error()
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		switch {
		case err.Error() == helpers.ErrWalletDisabled:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = map[string]interface{}{"transactions": res}
	json.NewEncoder(w).Encode(resp)
	return
}
