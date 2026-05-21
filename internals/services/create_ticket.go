package services

import (
	"context"
	"flight_booking_system/internals/repository/admin_repository"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GenerateTickets(pool *pgxpool.Pool, businessSeatsTotal, economySeatsTotal int, businessSeatsPrice, economySeatsPrice float64, flightID int, airplaneTailNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction %v ", err)
	}
	defer tx.Rollback(ctx)

	airplane, err := admin_repository.GetAirplaneByTailNumber(tx, ctx, airplaneTailNumber)
	if err != nil {
		return err
	}

	fmt.Printf("DEBUG: ReqBus=%d, ReqEco=%d, DBTotal=%d\n",
		businessSeatsTotal, economySeatsTotal, airplane.TotalSeats)

	if businessSeatsTotal+economySeatsTotal != airplane.TotalSeats {
		return fmt.Errorf("the sum of seats you enter does not sum to the number of seats in the airplane")
	}

	flight, err := admin_repository.GetFlightByID(tx, ctx, flightID, airplane.ID)
	if err != nil {
		return err
	}

	seats, err := admin_repository.GetSeatsByAirplaneID(tx, ctx, flight.AirplaneID)
	if err != nil {
		return err
	}

	for i, seat := range seats {
		if i < businessSeatsTotal {
			_, err = admin_repository.InsertTicket(tx, ctx, flightID, seat.ID, businessSeatsPrice, "business", "available")
			if err != nil {
				return err
			}
		} else {
			_, err = admin_repository.InsertTicket(tx, ctx, flightID, seat.ID, economySeatsPrice, "economy", "available")
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction %v ", err)
	}
	return nil
}
