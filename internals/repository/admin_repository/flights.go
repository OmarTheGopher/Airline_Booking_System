package admin_repository

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func InsertFlight(tx pgx.Tx, ctx context.Context, airplaneID int, depAirportID, arrAirportID int, depAt, arrAt time.Time) (*models.Flight, error) {

	query := `
		INSERT INTO flights(airplane_id, departure_airport_id, arrival_airport_id, departure_at, arrival_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, airplane_id, departure_airport_id, arrival_airport_id, departure_at, arrival_at, scheduled_at
	`

	var flight models.Flight
	err := tx.QueryRow(ctx, query, airplaneID, depAirportID, arrAirportID, depAt, arrAt).Scan(
		&flight.ID,
		&flight.AirplaneID,
		&flight.DepartureAirportID,
		&flight.ArrivalAirportID,
		&flight.DepartureAt,
		&flight.ArrivalAt,
		&flight.ScheduledAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "check_valid_journey") {
			return nil, fmt.Errorf("departure and arrival airports are the same")
		}
		if strings.Contains(err.Error(), "check_flight_duration") {
			return nil, fmt.Errorf("you are not light to travel back with time")
		}
		return nil, err
	}

	return &flight, nil
}

func GetFlightByID(tx pgx.Tx, ctx context.Context, flightID, airplaneID int) (*models.Flight, error) {

	query := `
		SELECT id, airplane_id, departure_airport_id, arrival_airport_id, departure_at, arrival_at, scheduled_at
		FROM flights	
		WHERE id = $1 AND airplane_id = $2;
	`
	var flight models.Flight
	err := tx.QueryRow(ctx, query, flightID, airplaneID).Scan(
		&flight.ID,
		&flight.AirplaneID,
		&flight.DepartureAirportID,
		&flight.ArrivalAirportID,
		&flight.DepartureAt,
		&flight.ArrivalAt,
		&flight.ScheduledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("flight not found")
		}
		return nil, err
	}

	return &flight, nil
}
