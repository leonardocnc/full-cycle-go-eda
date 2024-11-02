package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Balance struct {
	ID        string
	Amount    int
	AccountID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBalance(accountID string, amount int) (*Balance, error) {
	balance := &Balance{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Amount:    amount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := balance.Validate()
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Balance) Validate() error {
	if c.AccountID == "" {
		return errors.New("accountID is required")
	}
	return nil
}

func (b *Balance) Add(amount int) {
	b.Amount += amount
	b.UpdatedAt = time.Now()
}

func (b *Balance) Subtract(amount int) {
	b.Amount -= amount
	b.UpdatedAt = time.Now()
}
