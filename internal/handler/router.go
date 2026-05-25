package handler

import "github.com/gin-gonic/gin"

type TaskHandler interface {
	GetAllTasks(c *gin.Context)
	CreateTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	DeleteTaskByID(c *gin.Context)
	DeleteAllTasks(c *gin.Context)
	GetTasksByDue(c *gin.Context)
}

func RegisterRoutes(server *gin.Engine, taskHandler TaskHandler) {
	server.GET("/tasks", taskHandler.GetAllTasks)
	server.POST("/task", taskHandler.CreateTask)
	server.GET("/task/:id", taskHandler.GetTaskByID)
	server.DELETE("/task/:id", taskHandler.DeleteTaskByID)
	server.DELETE("/tasks", taskHandler.DeleteAllTasks)
	server.GET("/task/:year/:month/:day", taskHandler.GetTasksByDue)
}
