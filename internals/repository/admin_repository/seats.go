package admin_repository

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"

	"github.com/jackc/pgx/v5"
)

func GetSeatsByAirplaneID(tx pgx.Tx, ctx context.Context, airplaneID int) ([]models.Seat, error) {

	query := `
		SELECT id, airplane_id, seat_number
		FROM seats
		WHERE airplane_id = $1
		ORDER BY seat_number ASC;
	`

	rows, err := tx.Query(ctx, query, airplaneID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat

		err := rows.Scan(
			&seat.ID,
			&seat.AirplaneID,
			&seat.SeatNumber,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, pgx.ErrNoRows
			}
			return nil, err
		}

		seats = append(seats, seat)
	}

	return seats, nil
}
