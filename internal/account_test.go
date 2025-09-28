package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	t.Run("Deve criar uma conta de passageiro", func(t *testing.T) {
		account, err := CreateAccount("John Doe", "john.doe@gmail.com", "97456321558", "", "123456", true, false)
		require.NoError(t, err)
		require.NoError(t, err)
		assert.NotEmpty(t, account.AccountID)
		assert.Equal(t, "John Doe", account.Name)
		assert.Equal(t, "john.doe@gmail.com", account.Email)
		assert.Equal(t, "97456321558", account.CPF)
		assert.Equal(t, "123456", account.Password)
		assert.Equal(t, true, account.IsPassenger)
	})
	t.Run("Deve criar uma conta de motorista", func(t *testing.T) {
		account, err := CreateAccount("John Doe", "john.doe@gmail.com", "97456321558", "AAA9999", "123456", false, true)
		require.NoError(t, err)
		assert.NotEmpty(t, account.AccountID)
		assert.Equal(t, "John Doe", account.Name)
		assert.Equal(t, "john.doe@gmail.com", account.Email)
		assert.Equal(t, "97456321558", account.CPF)
		assert.Equal(t, "123456", account.Password)
		assert.Equal(t, "AAA9999", account.CarPlate)
		assert.Equal(t, true, account.IsDriver)
	})
	t.Run("Não deve criar uma conta de passageiro com o nome inválido", func(t *testing.T) {
		_, err := CreateAccount("John", "john.doe@gmail.com", "97456321558", "", "123456", true, false)
		assert.ErrorIs(t, err, ErrInvalidName)
	})
	t.Run("Não deve criar uma conta de passageiro com o email inválido", func(t *testing.T) {
		_, err := CreateAccount("John Doe", "john.doe", "97456321558", "", "123456", true, false)
		assert.ErrorIs(t, err, ErrInvalidEmail)
	})
	t.Run("Não deve criar uma conta de passageiro com cpf inválido", func(t *testing.T) {
		_, err := CreateAccount("John Doe", "john.doe@gmail.com", "9745632155", "", "123456", true, false)
		assert.ErrorIs(t, err, ErrInvalidCpf)
	})
	t.Run("Não deve criar uma conta de motorista com placa do carro inválida", func(t *testing.T) {
		_, err := CreateAccount("John Doe", "john.doe@gmail.com", "97456321558", "AAA999", "123456", false, true)
		assert.ErrorIs(t, err, ErrInvalidCarPlate)
	})
}
