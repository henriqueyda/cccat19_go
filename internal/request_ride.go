package internal

import (
	"context"
	"errors"
)

var (
	ErrOnlyPassengerCanRequestRide = errors.New("only passengers can request a ride")
	ErrRideInProgressForPassenger  = errors.New("there is already a ride in progress for this passenger")
)

type RequestRideData interface {
	SaveRide(ctx context.Context, ride Ride) (string, error)
	HasActiveRideByPassengerID(ctx context.Context, passengerID string) (bool, error)
}

type RequestRideInput struct {
	PassengerID string  `json:"passenger_id"`
	FromLat     float64 `json:"from_lat"`
	FromLong    float64 `json:"from_long"`
	ToLat       float64 `json:"to_lat"`
	ToLong      float64 `json:"to_long"`
}

type RequestRide struct {
	requestRideData RequestRideData
	getAccountData  GetAccountData
}

func NewRequestRide(requestRideData RequestRideData, getAccount GetAccountData) *RequestRide {
	return &RequestRide{
		requestRideData: requestRideData,
		getAccountData:  getAccount,
	}
}
func (r *RequestRide) Execute(ctx context.Context, input RequestRideInput) (string, error) {
	account, err := r.getAccountData.GetAccountByID(ctx, input.PassengerID)
	if err != nil {
		return "", err
	}
	if !account.IsPassenger {
		return "", ErrOnlyPassengerCanRequestRide
	}
	hasActiveRide, err := r.requestRideData.HasActiveRideByPassengerID(ctx, input.PassengerID)
	if err != nil {
		return "", err
	}
	if hasActiveRide {
		return "", ErrRideInProgressForPassenger
	}
	ride, err := CreateRide(input.PassengerID, input.FromLat, input.FromLong, input.ToLat, input.ToLong)
	if err != nil {
		return "", err
	}
	return r.requestRideData.SaveRide(ctx, ride)
}
