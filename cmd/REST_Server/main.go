package main

import (
	"REST_Server/internal/taskstore"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	router := gin.Default()
	taskstore.Open()
	defer taskstore.Close()

	router.POST("/task/", createTaskHandler)
	router.GET("/task/", getAllTasksHandler)
	router.DELETE("/task/", deleteAllTasksHandler)
	router.GET("/task/:id", getTaskHandler)
	router.DELETE("/task/:id", deleteTaskHandler)
	router.GET("/due/:year/:month/:day", dueHandler)

	router.Run()
}

func initEnv() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println("No .env file found")
	}
}

func createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if rt.Text == "" || rt.Due.IsZero() {
		c.String(http.StatusBadRequest, "Task or Due is required")
		return
	}

	id := taskstore.CreateTask(rt.Text, rt.Due)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getTaskHandler(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	task, err := taskstore.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": task.Id, "text": task.Text, "due": task.Due})
}

func getAllTasksHandler(c *gin.Context) {
	allTasks, err := taskstore.GetAllTasks()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, allTasks)
}

func deleteAllTasksHandler(c *gin.Context) {
	msg, err := taskstore.DeleteAllTasks()
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": msg})
}

func deleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		log.Println(err)
		return
	}

	msg, err := taskstore.DeleteTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.String(http.StatusOK, msg)
}

func dueHandler(c *gin.Context) {
	badRequestError := func() {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>, got %v", c.FullPath())
	}

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		badRequestError()
		return
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		badRequestError()
		return
	}

	tasks, err := taskstore.GetTasksByDueDate(year, time.Month(month), day)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
