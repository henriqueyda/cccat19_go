package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	// ctx := context.Background()
	t.Run("Deve criar uma conta de passageiro", func(t *testing.T) {
		cleanUpDB(t, "postgres://postgres:123456@localhost:5432/app")
		input := SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		}
		jsonInput, err := json.Marshal(input)
		assert.NoError(t, err)

		respSignup, err := http.Post("http://localhost:8080/signup", "application/json", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, respSignup.StatusCode)
		defer respSignup.Body.Close()
		body, err := io.ReadAll(respSignup.Body)
		assert.NoError(t, err)
		var output struct {
			AccountID string `json:"account_id"`
		}
		err = json.Unmarshal(body, &output)
		assert.NoError(t, err)

		respGet, err := http.Get(fmt.Sprintf("http://localhost:8080/accounts/%s", output.AccountID))
		assert.NoError(t, err)
		defer respGet.Body.Close()
		var outputGetAccount Account
		err = json.NewDecoder(respGet.Body).Decode(&outputGetAccount)
		assert.NoError(t, err)
		assert.NotEmpty(t, outputGetAccount.AccountID)
		assert.Equal(t, "John Doe", outputGetAccount.Name)
		assert.Equal(t, "john.doe@gmail.com", outputGetAccount.Email)
		assert.Equal(t, "97456321558", outputGetAccount.CPF)
		assert.Equal(t, "123456", outputGetAccount.Password)
		assert.Equal(t, true, outputGetAccount.IsPassenger)
	})

	t.Run("Não deve criar uma conta de passageiro com o nome inválido", func(t *testing.T) {
		cleanUpDB(t, "postgres://postgres:123456@localhost:5432/app")
		input := SignupInput{
			Name:        "John",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		}
		jsonInput, err := json.Marshal(input)
		assert.NoError(t, err)

		respSignup, err := http.Post("http://localhost:8080/signup", "application/json", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, respSignup.StatusCode)
		defer respSignup.Body.Close()
		body, err := io.ReadAll(respSignup.Body)
		assert.NoError(t, err)
		var apiError APIError
		err = json.Unmarshal(body, &apiError)
		assert.NoError(t, err)
		assert.Equal(t, "Error signing up: invalid name", apiError.Msg)
	})
}

func cleanUpDB(t *testing.T, connString string) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM ccca.account")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})
}
