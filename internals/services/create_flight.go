package services

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"
	"flight_booking_system/internals/repository/admin_repository"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateFlight(pool *pgxpool.Pool, tailNumber string, depAirportName, arrAirportName string, depTime, arrTime time.Time) (*models.Flight, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction")
	}

	defer tx.Rollback(ctx)

	airplane, err := admin_repository.GetAirplaneByTailNumber(tx, ctx, tailNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no airplane with that tail")
		}
		return nil, err
	}

	depAirport, err := admin_repository.GetAirportByName(tx, ctx, depAirportName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("departure airport not found")
		}
		return nil, err
	}
	arrAirport, err := admin_repository.GetAirportByName(tx, ctx, arrAirportName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("arrival airport not found")
		}
		return nil, err
	}

	flight, err := admin_repository.InsertFlight(tx, ctx, airplane.ID, depAirport.ID, arrAirport.ID, depTime, arrTime)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return flight, nil
}
