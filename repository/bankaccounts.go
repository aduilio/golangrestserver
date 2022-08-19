package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.comn/aduilio/golangrestserver/domain"
)

const create string = `
  CREATE TABLE IF NOT EXISTS bankaccounts (
  id INTEGER NOT NULL PRIMARY KEY,
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

func (b *BankAccountsDb) Save(account domain.Account) error {
	stmt, err := b.db.Prepare(`insert into bankaccounts(id, number, balance) values ($1, $2, $3)`)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	_, err = stmt.Exec(
		account.ID,
		account.Number,
		account.Balance,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func setupDb() *sql.DB {
	db, err := sql.Open("sqlite3", "db/bankaccounts.db")
	if err != nil {
		fmt.Println("Error connection to db", err.Error())
		return nil
	}
	if _, err := db.Exec(create); err != nil {
		return nil
	}

	return db
}
