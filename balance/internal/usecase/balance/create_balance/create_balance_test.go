package create_balance

import (
	"context"
	"testing"

	"balances/internal/entity"
	"balances/internal/event"
	"balances/internal/usecase/mocks"
	"balances/pkg/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type BalanceGatewayMock struct {
	mock.Mock
}

func (m *BalanceGatewayMock) Save(balance *entity.Balance) error {
	args := m.Called(balance)
	return args.Error(0)
}

func (m *BalanceGatewayMock) Get(id string) (*entity.Balance, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Balance), args.Error(1)
}

// func TestCreateBalanceUseCase_Execute(t *testing.T) {
// 	m := &BalanceGatewayMock{}
// 	m.On("Save", mock.Anything).Return(nil)

// 	uc := NewCreateBalanceUseCase(m)

// 	output, err := uc.Execute(CreateBalanceInputDTO{
// 		AccountID: "1",
// 		Amount:    100,
// 	})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, output)
// 	assert.Equal(t, "1", output.ID)
// 	assert.Equal(t, 100, output.Amount)
// 	m.AssertExpectations(t)
// 	m.AssertNumberOfCalls(t, "Save", 1)
// }

func TestCreateBalanceUseCase_Execute(t *testing.T) {
	balance, _ := entity.NewBalance("1", 100)

	// client2, _ := entity.NewClient("client2", "j@j2.com")
	// account2 := entity.NewAccount(client2)
	// account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateBalanceInputDTO{
		AccountID: balance.ID,
		Amount:    balance.Amount,
	}

	dispatcher := events.NewEventDispatcher()
	eventBalance := event.NewBalanceUpdated()
	ctx := context.Background()

	uc := NewCreateBalanceUseCase(mockUow, dispatcher, eventBalance)
	output, err := uc.Execute(ctx, inputDto)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
