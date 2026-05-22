package handler

import (
	"REST_Server/internal/model"
	"REST_Server/internal/taskstore"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskByID(ids string) (model.Task, error)
	CreateTask(description string, due time.Time) (int, error)
	DeleteTaskByID(ids string) (string, error)
	DeleteAllTasks() (string, error)
	GetTaskByDueDate(year string, month string, day string) ([]model.Task, error)
}

type handler struct {
	service TaskService
}

func NewHandler(service TaskService) *handler {
	return &handler{service}
}

func (h *handler) GetAllTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println(err) //?
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *handler) GetTaskByID(c *gin.Context) {
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

func (h *handler) createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if rt.Text == "" {
		c.String(http.StatusBadRequest, "Task is required")
		return
	}
	if rt.Due.IsZero() {
		c.String(http.StatusBadRequest, "Due is required")
		return
	}

	id, err := taskstore.CreateTask(rt.Text, rt.Due)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
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

func deleteAllTasksHandler(c *gin.Context) {
	msg, err := taskstore.DeleteAllTasks()
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": msg})
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
