package admin_repository

import (
	"context"
	"errors"
	"flight_booking_system/internals/models"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GenerateSeats(tx pgx.Tx, ctx context.Context, totalSeats int, airplaneID int) error {

	batch := &pgx.Batch{}

	query := `
        INSERT INTO seats(seat_number, airplane_id)
        VALUES($1, $2)
    `

	for i := 0; i < totalSeats; i++ {
		row := (i / 6) + 1
		letter := 'a' + rune(i%6)
		seatNumber := fmt.Sprintf("%d%c", row, letter)

		batch.Queue(query, seatNumber, airplaneID)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < totalSeats; i++ {
		_, err := br.Exec()
		if err != nil {
			// لو حصل مشكلة Unique (الكرسي متسجل قبل كدا)
			if strings.Contains(err.Error(), "unique") {
				return fmt.Errorf("seat already registered")
			}
			return err
		}
	}

	return nil
}

func AddAirplane(pool *pgxpool.Pool, model, tailNumber string, totalSeats int) (*models.Airplane, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	query := `
		INSERT INTO airplanes(model, tail_number, total_seats, is_active)
		VALUES($1, $2, $3, TRUE)
		RETURNING id, model, tail_number, total_seats, created_at, is_active
	`
	var airplane models.Airplane
	err = tx.QueryRow(ctx, query, model, tailNumber, totalSeats).Scan(
		&airplane.ID,
		&airplane.Model,
		&airplane.TailNumber,
		&airplane.TotalSeats,
		&airplane.CreatedAt,
		&airplane.IsActive,
	)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("aiplane already signed")
		}
		return nil, err
	}

	err = GenerateSeats(tx, ctx, airplane.TotalSeats, airplane.ID)

	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)
	return &airplane, nil
}

func GetAirplaneByTailNumber(tx pgx.Tx, ctx context.Context, tailNumber string) (*models.Airplane, error) {

	query := `
		SELECT id, model, tail_number, total_seats, created_at, is_active
		FROM airplanes
		WHERE tail_number = $1
	`

	var airplane models.Airplane
	err := tx.QueryRow(ctx, query, tailNumber).Scan(
		&airplane.ID,
		&airplane.Model,
		&airplane.TailNumber,
		&airplane.TotalSeats,
		&airplane.CreatedAt,
		&airplane.IsActive,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &airplane, nil
}
