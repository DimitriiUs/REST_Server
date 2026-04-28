package taskstore

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	err  error
)

func getDBUrl() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
}

func Open() {
	pool, err = pgxpool.New(context.Background(), getDBUrl())
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	log.Println("Connected!")
}

func Close() {
	pool.Close()
}
