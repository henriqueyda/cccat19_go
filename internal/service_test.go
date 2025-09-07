package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServices(t *testing.T) {
	ctx := context.Background()
	t.Run("Deve criar uma conta de passageiro", func(t *testing.T) {
		signup, getAccount := setupTest()
		signupID, err := signup.Signup(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		})
		require.NoError(t, err)
		outputGetAccount, err := getAccount.GetAccountByID(ctx, signupID)
		require.NoError(t, err)
		assert.NotEmpty(t, outputGetAccount.AccountID)
		assert.Equal(t, "John Doe", outputGetAccount.Name)
		assert.Equal(t, "john.doe@gmail.com", outputGetAccount.Email)
		assert.Equal(t, "97456321558", outputGetAccount.CPF)
		assert.Equal(t, "123456", outputGetAccount.Password)
		assert.Equal(t, true, outputGetAccount.IsPassenger)
	})
	t.Run("Deve criar uma conta de motorista", func(t *testing.T) {
		signup, getAccount := setupTest()
		signupID, err := signup.Signup(ctx, SignupInput{
			Name:     "John Doe",
			Email:    "john.doe@gmail.com",
			CPF:      "97456321558",
			Password: "123456",
			CarPlate: "AAA9999",
			IsDriver: true,
		})
		require.NoError(t, err)
		outputGetAccount, err := getAccount.GetAccountByID(ctx, signupID)
		require.NoError(t, err)
		assert.NotEmpty(t, outputGetAccount.AccountID)
		assert.Equal(t, "John Doe", outputGetAccount.Name)
		assert.Equal(t, "john.doe@gmail.com", outputGetAccount.Email)
		assert.Equal(t, "97456321558", outputGetAccount.CPF)
		assert.Equal(t, "123456", outputGetAccount.Password)
		assert.Equal(t, "AAA9999", outputGetAccount.CarPlate)
		assert.Equal(t, true, outputGetAccount.IsDriver)
	})
	t.Run("Não deve criar uma conta de passageiro com o nome inválido", func(t *testing.T) {
		signup, _ := setupTest()
		_, err := signup.Signup(ctx, SignupInput{
			Name:        "John",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		})
		assert.ErrorIs(t, err, ErrInvalidName)
	})
	t.Run("Não deve criar uma conta de passageiro com o email inválido", func(t *testing.T) {
		signup, _ := setupTest()
		_, err := signup.Signup(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		})
		assert.ErrorIs(t, err, ErrInvalidEmail)
	})
	t.Run("Não deve criar uma conta de passageiro com cpf inválido", func(t *testing.T) {
		signup, _ := setupTest()
		_, err := signup.Signup(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "9745632155",
			Password:    "123456",
			IsPassenger: true,
		})
		assert.ErrorIs(t, err, ErrInvalidCpf)
	})
	t.Run("Não deve criar uma conta de motorista com placa do carro inválida", func(t *testing.T) {
		signup, _ := setupTest()
		_, err := signup.Signup(ctx, SignupInput{
			Name:     "John Doe",
			Email:    "john.doe@gmail.com",
			CPF:      "97456321558",
			Password: "123456",
			CarPlate: "AAA999",
			IsDriver: true,
		})
		assert.ErrorIs(t, err, ErrInvalidCarPlate)
	})
	t.Run("Não deve criar uma conta de passageiro com conta duplicada", func(t *testing.T) {
		signup, _ := setupTest()
		input := SignupInput{
			Name:     "John Doe",
			Email:    "john.doe@gmail.com",
			CPF:      "97456321558",
			Password: "123456",
			CarPlate: "AAA9999",
			IsDriver: true,
		}
		_, err := signup.Signup(ctx, input)
		require.NoError(t, err)
		_, err = signup.Signup(ctx, input)
		assert.ErrorIs(t, err, ErrAccountAlreadyExists)
	})
}

func setupTest() (*Signup, *GetAccount) {
	accountDAO := NewAccountDAOMemory(&[]Account{})
	// accountDAO := NewAccountDAODatabase()
	mailerGateway := NewMailerGatewayMemory()
	signup := NewSignup(accountDAO, mailerGateway)
	getAccount := NewGetAccount(accountDAO)
	return signup, getAccount
}
