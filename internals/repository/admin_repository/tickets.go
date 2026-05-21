package admin_repository

import (
	"context"
	"flight_booking_system/internals/models"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func InsertTicket(tx pgx.Tx, ctx context.Context, flightID int, seatID int, price float64, class string, status string) (*models.Ticket, error) {

	query := `
		INSERT INTO tickets(flight_id, seat_id, price, class, status)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, flight_id, seat_id, price, class, status, created_at, updated_at
	`
	var ticket models.Ticket
	err := tx.QueryRow(ctx, query, flightID, seatID, price, class, status).Scan(
		&ticket.ID,
		&ticket.FlightID,
		&ticket.SeatID,
		&ticket.Price,
		&ticket.Class,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("ticket already created")
		}
		if strings.Contains(err.Error(), "check_valid_class") {
			return nil, fmt.Errorf("invalid class")
		}
		if strings.Contains(err.Error(), "check_valid_price") {
			return nil, fmt.Errorf("price cant be negative")
		}
		if strings.Contains(err.Error(), "check_valid_status") {
			return nil, fmt.Errorf("invalid status")
		}
		return nil, err
	}

	return &ticket, nil
}
