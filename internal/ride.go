package internal

import (
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	RideID      string    `json:"ride_id"`
	PassengerID string    `json:"passenger_id"`
	DriverID    *string   `json:"driver_id"`
	Status      string    `json:"status"`
	Fare        float64   `json:"fare"`
	Distance    float64   `json:"distance"`
	FromLat     float64   `json:"from_lat"`
	FromLong    float64   `json:"from_long"`
	ToLat       float64   `json:"to_lat"`
	ToLong      float64   `json:"to_long"`
	Date        time.Time `json:"date"`
}

func NewRide(rideID string, passengerID string, driverID *string, status string, fare, distance, fromLat, fromLong, toLat, toLong float64, date time.Time) (Ride, error) {
	return Ride{
		RideID:      rideID,
		PassengerID: passengerID,
		DriverID:    driverID,
		Status:      status,
		Fare:        fare,
		Distance:    distance,
		FromLat:     fromLat,
		FromLong:    fromLong,
		ToLat:       toLat,
		ToLong:      toLong,
		Date:        date,
	}, nil
}

func CreateRide(passengerID string, fromLat, fromLong, toLat, toLong float64) (Ride, error) {
	return NewRide(
		uuid.NewString(),
		passengerID,
		nil,
		"requested",
		0.0,
		0.0,
		fromLat,
		fromLong,
		toLat,
		toLong,
		time.Now().UTC(),
	)
}
