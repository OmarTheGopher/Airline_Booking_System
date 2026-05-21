package services

import (
	"context"
	"flight_booking_system/internals/models"
	"flight_booking_system/internals/repository"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BookTicket(pool *pgxpool.Pool, ticketID int, userID string) (*models.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction")
	}
	defer tx.Rollback(ctx)

	ticket, err := repository.GetTicketByID(tx, ctx, ticketID)
	if err != nil {
		return nil, err
	}

	ticket, err = repository.UpdateTicketStatus(tx, ctx, ticket.ID)
	if err != nil {
		return nil, err
	}

	booking, err := repository.InsertBooking(tx, ctx, ticket.ID, userID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return booking, nil
}
