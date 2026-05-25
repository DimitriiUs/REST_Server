package main

import (
	"REST_Server/internal/handler"
	"REST_Server/internal/repository/postgresql"
	"REST_Server/internal/service"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	pool, err := pgxpool.New(context.Background(), getDBUrl())
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer pool.Close()

	repo := postgresql.NewRepo(pool)
	taskService := service.NewService(repo)
	taskHandler := handler.NewHandler(taskService)

	server := gin.Default()
	handler.RegisterRoutes(server, taskHandler)

	server.Run()
}

func getDBUrl() string {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
}
