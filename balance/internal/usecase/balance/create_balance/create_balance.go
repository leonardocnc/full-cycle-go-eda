package create_balance

import (
	"balances/internal/entity"
	"balances/internal/gateway"
	"balances/pkg/events"
	"balances/pkg/uow"
	"context"
	"fmt"
	"time"
)

type CreateBalanceInputDTO struct {
	AccountID string `json:"account_id"`
	Amount    int    `json:"amount"`
}

type CreateBalanceOutputDTO struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBalanceUseCase struct {
	Uow            uow.UowInterface
	BalanceUpdated events.EventInterface
}

func NewCreateBalanceUseCase(
	Uow uow.UowInterface,
) *CreateBalanceUseCase {
	return &CreateBalanceUseCase{
		Uow: Uow,
	}
}

func (uc *CreateBalanceUseCase) Execute(ctx context.Context, input CreateBalanceInputDTO) (*CreateBalanceOutputDTO, error) {
	output := &CreateBalanceOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		balanceRepository := uc.getBalanceRepository(ctx)

		balance, err := entity.NewBalance(input.AccountID, input.Amount)
		if err != nil {
			return err
		}

		err = balanceRepository.Save(balance)

		if err != nil {
			return err
		}

		fmt.Printf("Balance: %+v\n", balance)

		output.ID = balance.ID
		output.AccountID = balance.AccountID
		output.Amount = balance.Amount
		output.CreatedAt = balance.CreatedAt

		return nil
	})

	if err != nil {
		fmt.Println("Error", err)
		return nil, err
	}

	return output, nil
}

func (uc *CreateBalanceUseCase) getBalanceRepository(ctx context.Context) gateway.BalanceGateway {
	repo, err := uc.Uow.GetRepository(ctx, "BalanceDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.BalanceGateway)
}
