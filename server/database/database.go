package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() *pgxpool.Pool {
	connectionString := "host=localhost port=5433 dbname=autorender-dev user=postgres password=admin sslmode=prefer connect_timeout=10"

	var err error
	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return Pool
}
