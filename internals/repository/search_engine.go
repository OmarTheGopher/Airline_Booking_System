package repository

import (
	"context"
	"flight_booking_system/internals/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SearchTicketEngine(pool *pgxpool.Pool, departureCity, arrivalCity string, departureTime, arrivalTime *time.Time, price *float64, class string) ([]models.SearchResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT tickets.id,
		tickets.price,
		flights.departure_at,
		flights.arrival_at,
		dep_airports.city,
		arr_airports.city
		FROM tickets
		JOIN flights ON tickets.flight_id = flights.id
		JOIN airports AS dep_airports ON flights.departure_airport_id = dep_airports.id  
		JOIN airports AS arr_airports ON flights.arrival_airport_id = arr_airports.id
		WHERE 
			($1::text IS NULL OR dep_airports.city ILIKE $1) AND
			($2::text IS NULL OR arr_airports.city ILIKE $2) AND
			($3::numeric IS NULL OR tickets.price <= $3) AND
			($4::date IS NULL OR flights.departure_at::DATE = $4) AND
			($5::date IS NULL OR flights.arrival_at::DATE = $5) AND
			($6::text IS NULL OR tickets.class = $6) AND
			(tickets.status = 'available')
		FOR UPDATE OF tickets
	`
	var cls interface{}
	if class == "" {
		cls = nil
	} else {
		cls = class
	}

	var dep interface{}
	if departureCity != "" {
		dep = "%" + departureCity + "%"
	} else {
		dep = nil
	}

	var arr interface{}
	if arrivalCity != "" {
		arr = "%" + arrivalCity + "%"
	} else {
		arr = nil
	}

	var results []models.SearchResult
	rows, err := pool.Query(ctx, query, dep, arr, price, departureTime, arrivalTime, cls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result models.SearchResult
		err := rows.Scan(
			&result.TicketID,
			&result.Price,
			&result.DepartureTime,
			&result.ArrivalTime,
			&result.DepartureCity,
			&result.ArrivalCity,
		)
		if err != nil {
			fmt.Printf("SCAN ERROR: %v\n", err) // اطبع الخطأ هنا فوراً
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
