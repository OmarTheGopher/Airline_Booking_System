package repository

import (
	"context"
	"flight_booking_system/internals/models"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func InsertBooking(tx pgx.Tx, ctx context.Context, ticketID int, userID string) (*models.Booking, error) {

	query := `
		INSERT INTO booking(ticket_id, user_id, status)
		VALUES($1, $2, 'confirmed')
		RETURNING id, user_id, ticket_id, status, created_at, updated_at
	`

	var booking models.Booking
	err := tx.QueryRow(ctx, query, ticketID, userID).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.TicketID,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("ticket already booked by this account")
		}
		return nil, err
	}
	return &booking, nil
}
