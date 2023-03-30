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

type AccountHandler struct {
	accService domain.AccountService
	timeout    time.Duration
}

func NewAccountHandler(acs domain.AccountService, timeout time.Duration) AccountHandler {
	handler := &AccountHandler{
		accService: acs,
		timeout:    timeout,
	}

	return *handler
}

func (o *AccountHandler) RegistUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err   error
		token = domain.Tokens{}
		resp  = helpers.Response{
			Status: helpers.SuccessMsg,
			Data:   nil,
		}
		customerXID = r.PostFormValue("customer_xid")
		errResp     = helpers.ErrResp{}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	if customerXID == "" {
		var Acc = domain.Accounts{CustomerXID: "Missing Required Field"}

		errResp.Err = Acc
		resp.Status = helpers.FailMsg
		resp.Data = errResp

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	token.Token, err = o.accService.Register(ctx, customerXID)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Data = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = token
	json.NewEncoder(w).Encode(resp)
	return
}
