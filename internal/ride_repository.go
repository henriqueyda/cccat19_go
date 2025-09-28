package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrRideNotFound = fmt.Errorf("ride not found")

type RideRepositoryDatabase struct {
}

func NewRideRepositoryDatabase() *RideRepositoryDatabase {
	return &RideRepositoryDatabase{}
}

func (r *RideRepositoryDatabase) SaveRide(ctx context.Context, ride Ride) (string, error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return "", err
	}
	_, err = conn.Exec(ctx, "insert into ccca.ride (ride_id, passenger_id, driver_id, status, fare, distance, from_lat, from_long, to_lat, to_long, date) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		&ride.RideID,
		&ride.PassengerID,
		&ride.DriverID,
		&ride.Status,
		&ride.Fare,
		&ride.Distance,
		&ride.FromLat,
		&ride.FromLong,
		&ride.ToLat,
		&ride.ToLong,
		&ride.Date,
	)
	if err != nil {
		return "", err
	}
	return ride.RideID, nil
}

func (a *RideRepositoryDatabase) GetRideByID(ctx context.Context, rideID string) (Ride, error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return Ride{}, err
	}
	var rRideID, rPassengerID, rStatus string
	var rDriverID *string
	var rFare, rDistance, rFromLat, rFromLong, rToLat, rToLong float64
	var rDate time.Time
	err = conn.QueryRow(ctx, "select ride_id, passenger_id, driver_id, status, fare, distance, from_lat, from_long, to_lat, to_long, date from ccca.ride where ride_id = $1", rideID).Scan(
		&rRideID,
		&rPassengerID,
		&rDriverID,
		&rStatus,
		&rFare,
		&rDistance,
		&rFromLat,
		&rFromLong,
		&rToLat,
		&rToLong,
		&rDate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Ride{}, ErrRideNotFound
		}
		return Ride{}, err
	}
	ride, err := NewRide(rRideID, rPassengerID, rDriverID, rStatus, rFare, rDistance, rFromLat, rFromLong, rToLat, rToLong, rDate)
	if err != nil {
		return Ride{}, err
	}
	return ride, nil
}

func (a *RideRepositoryDatabase) HasActiveRideByPassengerID(ctx context.Context, passengerID string) (bool, error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:123456@localhost:5432/app")
	defer func() {
		closeErr := conn.Close(ctx)
		if err != nil {
			if closeErr != nil {
				fmt.Printf("failed to close connection: %s", closeErr)
			}
			return
		}
		err = closeErr
	}()
	if err != nil {
		return false, err
	}
	var count int
	err = conn.QueryRow(ctx, "select count(*)::int as count from ccca.ride where passenger_id = $1 and status <> 'completed'", passengerID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
