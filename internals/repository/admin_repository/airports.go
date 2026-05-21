package admin_repository

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AddAirport(pool *pgxpool.Pool, name, city, country string) (*models.Airport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO airports(name, city, country, is_active)
		VALUES($1, $2, $3, TRUE)
		RETURNING id, name, city, country, created_at, is_active
	`

	var airport models.Airport
	err := pool.QueryRow(ctx, query, name, city, country).Scan(
		&airport.ID,
		&airport.Name,
		&airport.City,
		&airport.Country,
		&airport.CreatedAt,
		&airport.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &airport, nil
}

func GetAirportByName(tx pgx.Tx, ctx context.Context, name string) (*models.Airport, error) {
	query := `
		SELECT id, name, city, country, created_at, is_active
		FROM airports
		WHERE name = $1
	`

	var airport models.Airport
	err := tx.QueryRow(ctx, query, name).Scan(
		&airport.ID,
		&airport.Name,
		&airport.City,
		&airport.Country,
		&airport.CreatedAt,
		&airport.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err //Internal Server
	}

	return &airport, nil
}
