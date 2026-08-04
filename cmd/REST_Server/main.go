package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"REST_Server/internal/config"
	"REST_Server/internal/handler"
	"REST_Server/internal/logger"
	"REST_Server/internal/repository/postgresql"
	"REST_Server/internal/service"
)

func main() {
	cfg := config.GetConfig()
	lg,err:=logger.New(cfg.LogLevel)
	if err!=nil{
		lg.Error(err.Error())
	}
	

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.DBURL)
	if err != nil {
		lg.Error("%w: unable to connection to database",err.Error())
	}
	defer pool.Close()

	repo := postgresql.NewRepo(pool)
	taskService := service.NewService(repo,lg)
	taskHandler := handler.NewHandler(taskService)

	server := gin.Default()
	handler.RegisterRoutes(server, taskHandler)

	err = server.Run()
	if err != nil {
		log.Panicf("Unable to start server: %v\n", err)
	}
}
