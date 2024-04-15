package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func StartDBConnection() *pgx.Conn {

	// you must set url of your database. syntax allowed :
	// postgresql://user:password@localhost:5432/database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment is non set")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
