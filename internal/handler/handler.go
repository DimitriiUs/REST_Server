package handler

import (
	"REST_Server/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskByID(ids string) (model.Task, error)
	CreateTask(description string, due time.Time) (int, error)
	DeleteTaskByID(ids string) (string, error)
	DeleteAllTasks() (string, error)
	GetTasksByDue(year string, month string, day string) ([]model.Task, error)
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

	//Добавить ошибку ErrorNotFound, добавить проверку на ошибки
	task, err := h.service.GetTaskByID(ids)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": task.Id, "text": task.Text, "due": task.Due})
}

func (h *handler) CreateTask(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateTask(rt.Text, rt.Due)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *handler) DeleteTaskByID(c *gin.Context) {
	ids := c.Param("id")

	msg, err := h.service.DeleteTaskByID(ids)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.String(http.StatusNoContent, msg)
}

func (h *handler) DeleteAllTasks(c *gin.Context) {
	msg, err := h.service.DeleteAllTasks()
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		log.Println(err)
		return
	}
	c.String(http.StatusNoContent, msg)
}

func (h *handler) GetTasksByDue(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	day := c.Param("day")

	tasks, err := h.service.GetTasksByDue(year, month, day)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}
