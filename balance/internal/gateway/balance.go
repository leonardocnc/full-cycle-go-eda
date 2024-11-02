package gateway

import "balances/internal/entity"

type BalanceGateway interface {
	GetBalanceByAccountID(id string) (*entity.Balance, error)
	UpdateBalanceByAccountID(accountID string, amount int) error
	Save(balance *entity.Balance) error
}
