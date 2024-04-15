package main

import (
	"context"
	"fmt"
	"github.com/Sem4kok/restful-api/db"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func main() {

	// you must set url of your database. syntax allowed :
	// postgresql://user:password@localhost:5432/database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment in son set")
	}

	// connecting to PostgreSQL database storage
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()

	albums, err := db.GetAlbumsFromDB(context.Background(), conn)
	if err != nil && albums == nil {
		log.Fatal(err)
	} else if err != nil {
		fmt.Println("error, with rows.Err(): ", err.Error())
	}

}
