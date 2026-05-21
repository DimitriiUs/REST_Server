package main

import (
	"REST_Server/internal/handler"
	"REST_Server/internal/taskstore"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	taskstore.Open()
	defer taskstore.Close()

	router.POST("/task/", handler.createTaskHandler())
	router.GET("/task/", getAllTasksHandler)
	router.DELETE("/task/", deleteAllTasksHandler)
	router.GET("/task/:id", getTaskHandler)
	router.DELETE("/task/:id", deleteTaskHandler)
	router.GET("/due/:year/:month/:day", dueHandler)

	router.Run()
}
