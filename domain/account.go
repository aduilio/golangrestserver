package domain

import (
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	ID      string  `json:"id"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
}

func NewAccount() *Account {
	account := &Account{}
	account.ID = uuid.NewV4().String()
	account.Balance = 1000
	return account
}
