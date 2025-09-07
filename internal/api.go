package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var accounts = []Account{}

func JSONError(w http.ResponseWriter, error any, code int) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(error)
}

type APIError struct {
	Msg string `json:"msg"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// accountDAO := NewAccountDAOMemory(&accounts)
	accountDAO := NewAccountDAODatabase()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error reading request body: %s", err)}, http.StatusInternalServerError)
		return
	}
	var input SignupInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error reading request body: %s", err)}, http.StatusInternalServerError)
		return
	}

	mailerGateway := NewMailerGatewayMemory()

	signup := NewSignup(accountDAO, mailerGateway)
	accountID, err := signup.Signup(ctx, input)
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error signing up: %s", err)}, http.StatusUnprocessableEntity)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		AccountID string `json:"account_id"`
	}{
		AccountID: accountID,
	})
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error encoding response: %s", err)}, http.StatusInternalServerError)
		return
	}
}

func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// accountDAO := NewAccountDAOMemory(&accounts)
	accountDAO := NewAccountDAODatabase()
	getAccount := NewGetAccount(accountDAO)
	account, err := getAccount.GetAccountByID(ctx, r.PathValue("account_id"))
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error getting account: %s", err)}, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		JSONError(w, APIError{Msg: fmt.Sprintf("Error encoding response: %s", err)}, http.StatusInternalServerError)
		return
	}
}
