package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.comn/aduilio/golangrestserver/domain"
	"github.comn/aduilio/golangrestserver/dto"
	"github.comn/aduilio/golangrestserver/repository"
)

var dbRepository repository.BankAccountsDb

func main() {
	dbRepository = *repository.NewBankAccountsDb()

	router := mux.NewRouter()

	router.HandleFunc("/bank-accounts", PostBankAccounts).Methods("POST")
	router.HandleFunc("/bank-accounts/transfer", PostTranfer).Methods("POST")

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Server runnning at 8000")
}

func PostBankAccounts(w http.ResponseWriter, r *http.Request) {
	accountRequest, isValid := parseAccountRequest(w, r)
	if !isValid {
		return
	}

	if !dbRepository.ValidateNumber(accountRequest.Number) {
		createErrorMessage(w, "Fail to save the account", "This account number already exists")
		return
	}

	account := domain.NewAccount()
	account.Number = accountRequest.Number

	err := dbRepository.Save(account)
	if err != nil {
		createErrorMessage(w, "Fail to save the account", err.Error())
		return
	}

	response := dto.AccountResponse{ID: account.ID, Number: account.Number}
	createResponse(w, http.StatusCreated, response)
}

func PostTranfer(w http.ResponseWriter, r *http.Request) {
	transferRequest, isValid := parseTransferRequest(r, w)
	if !isValid {
		return
	}

	if dbRepository.ValidateNumber(transferRequest.From) {
		createErrorMessage(w, "Fail to transfer", "The source account number does not exist")
		return
	}

	if dbRepository.ValidateNumber(transferRequest.To) {
		createErrorMessage(w, "Fail to transfer", "The destination account number does not exist")
		return
	}

	accountFromBalance, accountToBalance, err := dbRepository.Tranfer(transferRequest)
	if err != nil {
		createErrorMessage(w, "Fail to transfer", err.Error())
		return
	}

	accountFrom := dto.AccountTransfer{Number: transferRequest.From, Balance: *accountFromBalance}
	accountTo := dto.AccountTransfer{Number: transferRequest.To, Balance: *accountToBalance}
	response := dto.TransferResponse{From: accountFrom, To: accountTo}

	createResponse(w, http.StatusOK, response)
}

func parseAccountRequest(w http.ResponseWriter, r *http.Request) (*dto.AccountRequest, bool) {
	accountRequest := dto.AccountRequest{}

	err := json.NewDecoder(r.Body).Decode(&accountRequest)
	if err != nil {
		createErrorMessage(w, "Fail to parse the request body", err.Error())
		return nil, false
	}

	if len(accountRequest.Number) == 0 {
		createErrorMessage(w, "Missing account_number", "Inform the account_number in the body")
		return nil, false
	}

	return &accountRequest, true
}

func parseTransferRequest(r *http.Request, w http.ResponseWriter) (*dto.TransferRequest, bool) {
	transferRequest := dto.TransferRequest{}

	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		createErrorMessage(w, "Fail to parse the request body", err.Error())
		return nil, false
	}

	if len(transferRequest.From) == 0 {
		createErrorMessage(w, "Missing from", "Inform the source account (from) in the body")
		return nil, false
	}

	if len(transferRequest.To) == 0 {
		createErrorMessage(w, "Missing to", "Inform the destination account (to) in the body")
		return nil, false
	}

	if transferRequest.Amount <= 0 {
		createErrorMessage(w, "Invalid amount", "Inform a positive amount")
		return nil, false
	}

	return &transferRequest, true
}

func createErrorMessage(w http.ResponseWriter, message string, details string) {
	error := dto.ErrorMessage{Message: message, Details: details}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(error)
}

func createResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}
