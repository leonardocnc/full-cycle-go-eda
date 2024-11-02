package database

import (
	"balances/internal/entity"
	"database/sql"
	"fmt"
)

type BalanceDB struct {
	DB *sql.DB
}

func NewBalanceDB(db *sql.DB) *BalanceDB {
	return &BalanceDB{
		DB: db,
	}
}

func (b *BalanceDB) GetBalanceByAccountID(id string) (*entity.Balance, error) {
	balance := &entity.Balance{}
	fmt.Println("id", id)
	stmt, err := b.DB.Prepare("SELECT id, amount, account_id, created_at FROM balances WHERE account_id = ? ORDER BY created_at DESC LIMIT 1;")
	if err != nil {
		fmt.Println("GetBalanceByAccountID - b.DB.Prepar")
		fmt.Println("Error", err)
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err := row.Scan(&balance.ID, &balance.Amount, &balance.AccountID, &balance.CreatedAt); err != nil {
		fmt.Println("GetBalanceByAccountID - row.Scan")
		fmt.Println("Error", err)
		return nil, err
	}
	return balance, nil
}

func (b *BalanceDB) UpdateBalanceByAccountID(accountID string, amount int) error {
	query := "UPDATE balances SET amount = amount + ? WHERE account_id = ?;"
	_, err := b.DB.Exec(query, amount, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (b *BalanceDB) Save(balance *entity.Balance) error {
	stmt, err := b.DB.Prepare("INSERT INTO balances (id, account_id, amount, created_at) VALUES (?, ?, ?, ?);")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(balance.ID, balance.AccountID, balance.Amount, balance.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
