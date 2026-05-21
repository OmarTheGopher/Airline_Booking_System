package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Println("Failed to parse database url")
		return nil, err
	}

	var ctx context.Context = context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Println("Failed to open connection pool")
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Println("Failed to ping database")
		pool.Close()
		return nil, err
	}

	log.Println("Connected succesfully!!")
	return pool, nil
}
