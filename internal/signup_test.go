package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignup(t *testing.T) {
	ctx := context.Background()
	t.Run("Não deve criar uma conta de passageiro com conta duplicada", func(t *testing.T) {
		signup, _ := setupSignupTest(t)
		input := SignupInput{
			Name:     "John Doe",
			Email:    "john.doe@gmail.com",
			CPF:      "97456321558",
			Password: "123456",
			CarPlate: "AAA9999",
			IsDriver: true,
		}
		_, err := signup.Execute(ctx, input)
		require.NoError(t, err)
		_, err = signup.Execute(ctx, input)
		assert.ErrorIs(t, err, ErrAccountAlreadyExists)
	})
}

func setupSignupTest(t *testing.T) (*Signup, *GetAccount) {
	accountRepository := NewAccountRepositoryMemory(&[]Account{})
	// accountRepository := NewAccountRepositoryDatabase()
	// cleanUpDB(t, "postgres://postgres:123456@localhost:5432/app")
	mailerGateway := NewMailerGatewayMemory()
	signup := NewSignup(accountRepository, mailerGateway)
	getAccount := NewGetAccount(accountRepository)
	return signup, getAccount
}
