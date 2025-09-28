package internal

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestRide(t *testing.T) {
	ctx := context.Background()
	t.Run("Deve criar uma solicitação de corrida se a conta for de passageiro", func(t *testing.T) {
		signup, requestRide, getRide := setupRequestRideTest(t)
		signupID, err := signup.Execute(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		})
		require.NoError(t, err)
		rideID, err := requestRide.Execute(ctx, RequestRideInput{
			PassengerID: signupID,
			FromLat:     25.1,
			FromLong:    30.1,
			ToLat:       35.1,
			ToLong:      40.1,
		})
		require.NoError(t, err)
		ride, err := getRide.Execute(ctx, rideID)
		require.NoError(t, err)
		assert.Equal(t, ride.RideID, rideID)
		assert.Equal(t, signupID, ride.PassengerID)
		assert.Equal(t, "John Doe", ride.PassagerName)
		assert.Equal(t, 25.1, ride.FromLat)
		assert.Equal(t, 30.1, ride.FromLong)
		assert.Equal(t, 35.1, ride.ToLat)
		assert.Equal(t, 40.1, ride.ToLong)
		assert.Equal(t, "requested", ride.Status)
		assert.Equal(t, 0.0, ride.Fare)
		assert.Equal(t, 0.0, ride.Distance)
		assert.WithinDuration(t, time.Now().UTC(), ride.Date, time.Second*5)
	})
	t.Run("Não deve criar uma solicitação de corrida se a conta não for de passageiro", func(t *testing.T) {
		signup, requestRide, _ := setupRequestRideTest(t)
		signupID, err := signup.Execute(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: false,
		})
		require.NoError(t, err)
		_, err = requestRide.Execute(ctx, RequestRideInput{
			PassengerID: signupID,
			FromLat:     25.1,
			FromLong:    30.1,
			ToLat:       35.1,
			ToLong:      40.1,
		})
		require.ErrorIs(t, err, ErrOnlyPassengerCanRequestRide)
	})
	t.Run("Não deve criar uma solicitação de corrida se a conta não existir", func(t *testing.T) {
		_, requestRide, _ := setupRequestRideTest(t)
		_, err := requestRide.Execute(ctx, RequestRideInput{
			PassengerID: uuid.NewString(),
			FromLat:     25.1,
			FromLong:    30.1,
			ToLat:       35.1,
			ToLong:      40.1,
		})
		require.ErrorIs(t, err, ErrAccountNotFound)
	})
	t.Run("Não deve criar uma solicitação de corrida se já existir uma solicitação em andamento para o passageiro", func(t *testing.T) {
		signup, requestRide, _ := setupRequestRideTest(t)
		signupID, err := signup.Execute(ctx, SignupInput{
			Name:        "John Doe",
			Email:       "john.doe@gmail.com",
			CPF:         "97456321558",
			Password:    "123456",
			IsPassenger: true,
		})
		require.NoError(t, err)
		_, err = requestRide.Execute(ctx, RequestRideInput{
			PassengerID: signupID,
			FromLat:     25.1,
			FromLong:    30.1,
			ToLat:       35.1,
			ToLong:      40.1,
		})
		require.NoError(t, err)
		_, err = requestRide.Execute(ctx, RequestRideInput{
			PassengerID: signupID,
			FromLat:     25.1,
			FromLong:    30.1,
			ToLat:       35.1,
			ToLong:      40.1,
		})
		require.ErrorIs(t, err, ErrRideInProgressForPassenger)
	})
}

func setupRequestRideTest(t *testing.T) (*Signup, *RequestRide, *GetRide) {
	cleanUpDB(t, "postgres://postgres:123456@localhost:5432/app")
	rideRepository := NewRideRepositoryDatabase()
	accountRepository := NewAccountRepositoryDatabase()
	mailerGateway := NewMailerGatewayMemory()
	signup := NewSignup(accountRepository, mailerGateway)
	getRide := NewGetRide(rideRepository, accountRepository)
	requestRide := NewRequestRide(rideRepository, accountRepository)
	return signup, requestRide, getRide
}
