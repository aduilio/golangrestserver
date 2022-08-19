package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.comn/aduilio/golangrestserver/dto"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/bank-accounts", PostBankAccounts).Methods("POST")
	router.HandleFunc("/bank-accounts/transfer", PostTranfer).Methods("POST")

	fmt.Println("Server runnning at 8080")
	http.ListenAndServe(":8000", router)
}

func PostBankAccounts(w http.ResponseWriter, r *http.Request) {
	accountRequest := dto.AccountRequest{}

	err := json.NewDecoder(r.Body).Decode(&accountRequest)
	if err != nil {
		createErrorMessage(w, "Fail to parse the request body", err.Error())
		return
	}

	if len(accountRequest.Number) == 0 {
		createErrorMessage(w, "Missing account_number", "Inform the account_number in the body")
		return
	}

	fmt.Println(fmt.Sprintf("Creating a new bank account: %s", accountRequest.Number))

	response := dto.AccountResponse{ID: "123456", Number: accountRequest.Number}

	createResponse(w, http.StatusCreated, response)
}

func PostTranfer(w http.ResponseWriter, r *http.Request) {
	transferRequest := dto.TransferRequest{}

	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		createErrorMessage(w, "Fail to parse the request body", err.Error())
		return
	}

	if len(transferRequest.From) == 0 {
		createErrorMessage(w, "Missing from", "Inform the source account (from) in the body")
		return
	}

	if len(transferRequest.To) == 0 {
		createErrorMessage(w, "Missing to", "Inform the destination account (to) in the body")
		return
	}

	if transferRequest.Amount <= 0 {
		createErrorMessage(w, "Invalid amount", "Inform a positive amount")
		return
	}

	fmt.Println(fmt.Sprintf("Transfering %f from %s to %s", transferRequest.Amount, transferRequest.From, transferRequest.To))

	accountFrom := dto.AccountTransfer{Number: transferRequest.From, Balance: 100}
	accountTo := dto.AccountTransfer{Number: transferRequest.To, Balance: 50}
	response := dto.TransferResponse{From: accountFrom, To: accountTo}

	createResponse(w, http.StatusOK, response)
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
