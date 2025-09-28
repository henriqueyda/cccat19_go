package internal

import (
	"context"
	"time"
)

// Here we can define the passenger's name, which is not stored in the Ride database
// That's why we create a new struct for this use case
type GetRideOutput struct {
	RideID       string    `json:"ride_id"`
	PassengerID  string    `json:"passenger_id"`
	PassagerName string    `json:"passenger_name"`
	DriverID     *string   `json:"driver_id"`
	FromLat      float64   `json:"from_lat"`
	FromLong     float64   `json:"from_long"`
	ToLat        float64   `json:"to_lat"`
	ToLong       float64   `json:"to_long"`
	Status       string    `json:"status"`
	Fare         float64   `json:"fare"`
	Distance     float64   `json:"distance"`
	Date         time.Time `json:"date"`
}

type GetRideData interface {
	GetRideByID(ctx context.Context, rideID string) (Ride, error)
}

type GetRide struct {
	getRideData    GetRideData
	getAccountData GetAccountData
}

func NewGetRide(getRideData GetRideData, getAccountData GetAccountData) *GetRide {
	return &GetRide{
		getRideData:    getRideData,
		getAccountData: getAccountData,
	}
}

func (g *GetRide) Execute(ctx context.Context, rideID string) (GetRideOutput, error) {
	ride, err := g.getRideData.GetRideByID(ctx, rideID)
	if err != nil {
		return GetRideOutput{}, err
	}
	passenger, err := g.getAccountData.GetAccountByID(ctx, ride.PassengerID)
	if err != nil {
		return GetRideOutput{}, err
	}
	return GetRideOutput{
		RideID:       ride.RideID,
		PassengerID:  ride.PassengerID,
		PassagerName: passenger.Name,
		DriverID:     ride.DriverID,
		FromLat:      ride.FromLat,
		FromLong:     ride.FromLong,
		ToLat:        ride.ToLat,
		ToLong:       ride.ToLong,
		Status:       ride.Status,
		Fare:         ride.Fare,
		Distance:     ride.Distance,
		Date:         ride.Date,
	}, nil
}
