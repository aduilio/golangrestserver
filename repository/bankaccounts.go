package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.comn/aduilio/golangrestserver/domain"
	"github.comn/aduilio/golangrestserver/dto"
)

const create string = `
  CREATE TABLE IF NOT EXISTS accounts (
  id TEXT NOT NULL PRIMARY KEY,
  number TEXT NOT NULL,
  balance REAL
  );`

type BankAccountsDb struct {
	db *sql.DB
}

func NewBankAccountsDb() *BankAccountsDb {
	db := setupDb()
	return &BankAccountsDb{db: db}
}

func (b *BankAccountsDb) ValidateNumber(accountNumber string) bool {
	var accountId string
	err := b.db.QueryRow("SELECT id FROM accounts WHERE number = ?", accountNumber).Scan(&accountId)
	if err != nil {
		return true
	}
	if len(accountId) != 0 {
		return false
	}

	return true
}

func (b *BankAccountsDb) Save(account *domain.Account) error {
	stmt, err := b.db.Prepare(`INSERT INTO accounts(id, number, balance) VALUES ($1, $2, $3)`)
	if err != nil {
		fmt.Println("Erro preparing the statement", err.Error())
		return err
	}
	_, err = stmt.Exec(
		account.ID,
		account.Number,
		account.Balance,
	)
	if err != nil {
		fmt.Println("Error executing the statement: ", err.Error())
		return err
	}
	err = stmt.Close()
	if err != nil {
		fmt.Println("Error closing the connection: ", err.Error())
		return err
	}
	return nil
}

func (b *BankAccountsDb) Tranfer(request *dto.TransferRequest) (*float64, *float64, error) {
	accountFromBalance := b.getBalance(request.From)

	if accountFromBalance < request.Amount {
		return nil, nil, errors.New("The source account does not have enough balance")
	}

	accountToBalance := b.getBalance(request.To)

	accountFromBalance -= request.Amount
	accountToBalance += request.Amount

	err := b.updateBalance(request.From, accountFromBalance)
	if err != nil {
		fmt.Println("Error updating the balance: ", err.Error())
		return nil, nil, err
	}
	err = b.updateBalance(request.To, accountToBalance)
	if err != nil {
		fmt.Println("Error updating the balance: ", err.Error())
		return nil, nil, err
	}

	return &accountFromBalance, &accountToBalance, nil
}

func setupDb() *sql.DB {
	db, err := sql.Open("sqlite3", "db/bank.db")
	if err != nil {
		fmt.Println("Error connection to db: ", err.Error())
		return nil
	}
	if _, err := db.Exec(create); err != nil {
		return nil
	}

	return db
}

func (b *BankAccountsDb) getBalance(accountNumber string) float64 {
	var balance float64
	err := b.db.QueryRow("SELECT balance FROM accounts WHERE number = ?", accountNumber).Scan(&balance)
	if err != nil {
		return 0
	}

	return balance
}

func (b *BankAccountsDb) updateBalance(accountNumber string, balance float64) error {
	stmt, err := b.db.Prepare(`UPDATE accounts SET balance = $1 WHERE number = $2`)
	if err != nil {
		fmt.Println("Erro preparing the statement", err.Error())
		return err
	}
	_, err = stmt.Exec(
		balance,
		accountNumber,
	)
	if err != nil {
		fmt.Println("Error executing the statement: ", err.Error())
		return err
	}
	err = stmt.Close()
	if err != nil {
		fmt.Println("Error closing the connection: ", err.Error())
		return err
	}
	return nil
}
