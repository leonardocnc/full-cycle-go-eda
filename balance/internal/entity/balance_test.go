package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	balance, err := NewBalance("1", 100)
	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, "1", balance.AccountID)
	assert.Equal(t, 100, balance.Amount)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	balance, err := NewBalance("", 0)
	assert.NotNil(t, err)
	assert.Nil(t, balance)
}

func TestUpdateBalacneWithInvalidArgs(t *testing.T) {
	balance, err := NewBalance("", 100)
	assert.Error(t, err, "accountID is required")
	assert.Nil(t, balance)
}

func TestAddBalance(t *testing.T) {
	balance, _ := NewBalance("1", 100)
	balance.Add(100)
	assert.Equal(t, 200, balance.Amount)
}

func TestSubtractBalance(t *testing.T) {
	balance, _ := NewBalance("1", 100)
	balance.Subtract(100)
	assert.Equal(t, 0, balance.Amount)
}
