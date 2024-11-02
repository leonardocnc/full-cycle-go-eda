package get_balance

import (
	"balances/internal/gateway"
	"fmt"
	"time"
)

type GetBalanceInputDTO struct {
	AccountID string `json:"account_id"`
}

type GetBalanceOutputDTO struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type GetBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewGetBalanceUseCase(b gateway.BalanceGateway) *GetBalanceUseCase {
	return &GetBalanceUseCase{
		BalanceGateway: b,
	}
}

func (g *GetBalanceUseCase) Execute(input GetBalanceInputDTO) (*GetBalanceOutputDTO, error) {
	fmt.Println("input.AccountID", input.AccountID)
	balance, err := g.BalanceGateway.GetBalanceByAccountID(input.AccountID)
	if err != nil {
		fmt.Println("Execute GetBalanceByAccountID")
		fmt.Println("Error", err)
		return &GetBalanceOutputDTO{}, err
	}
	fmt.Println("balance", balance)
	return &GetBalanceOutputDTO{
		ID:        balance.ID,
		AccountID: balance.AccountID,
		Amount:    balance.Amount,
		CreatedAt: balance.CreatedAt,
	}, nil
}
