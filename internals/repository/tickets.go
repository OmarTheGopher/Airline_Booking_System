package repository

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func UpdateTicketStatus(tx pgx.Tx, ctx context.Context, ticketID int) (*models.Ticket, error) {

	query := `
		UPDATE tickets
		SET status = 'booked', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, flight_id, seat_id, class, price, status, created_at, updated_at
	`

	var ticket models.Ticket
	err := tx.QueryRow(ctx, query, ticketID).Scan(
		&ticket.ID,
		&ticket.FlightID,
		&ticket.SeatID,
		&ticket.Class,
		&ticket.Price,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func GetTicketByID(tx pgx.Tx, ctx context.Context, ticketID int) (*models.Ticket, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, flight_id, seat_id, class, price, status, created_at, updated_at
		FROM tickets
		WHERE id = $1
		FOR UPDATE
	`

	var ticket models.Ticket
	err := tx.QueryRow(ctx, query, ticketID).Scan(
		&ticket.ID,
		&ticket.FlightID,
		&ticket.SeatID,
		&ticket.Class,
		&ticket.Price,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ticket not found")
		}
		return nil, err
	}
	return &ticket, nil
}
