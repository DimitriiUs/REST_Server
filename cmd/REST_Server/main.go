package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"REST_Server/internal/config"
	"REST_Server/internal/handler"
	"REST_Server/internal/repository/postgresql"
	"REST_Server/internal/service"
)

func main() {
	config.LoadPostgresConfig()

	pool, err := pgxpool.New(context.Background(), config.GetDBUrl())
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
