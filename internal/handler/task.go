package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	errs "REST_Server/internal/errors"
	"REST_Server/internal/model"
)

//go:generate go tool counterfeiter -o /fake . TaskService

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskByID(ids string) (model.Task, error)
	CreateTask(description string, due time.Time) (int, error)
	DeleteTaskByID(ids string) error
	DeleteAllTasks() error
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
		handleError(c, err)
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *handler) GetTaskByID(c *gin.Context) {
	ids := c.Param("id")

	task, err := h.service.GetTaskByID(ids)

	if err != nil {
		handleError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"id": task.ID, "text": task.Text, "due": task.Due})
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
		handleError(c, err)
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *handler) DeleteTaskByID(c *gin.Context) {
	ids := c.Param("id")

	err := h.service.DeleteTaskByID(ids)

	if err != nil {
		handleError(c, err)
	}

	c.String(http.StatusNoContent, "task was deleted")
}

func (h *handler) DeleteAllTasks(c *gin.Context) {
	err := h.service.DeleteAllTasks()

	if err != nil {
		handleError(c, err)
	}

	c.String(http.StatusNoContent, "all tasks was deleted")
}

func (h *handler) GetTasksByDue(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	day := c.Param("day")

	tasks, err := h.service.GetTasksByDue(year, month, day)

	if err != nil {
		handleError(c, err)
	}

	c.JSON(http.StatusOK, tasks)
}

func handleError(c *gin.Context,err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound):
		c.String(http.StatusNotFound, err.Error())
	case errors.Is(err, errs.ErrInvalidID):
		c.String(http.StatusBadRequest, err.Error())
	case errors.Is(err, errs.ErrInvalidDescription):
		c.String(http.StatusBadRequest, err.Error())
		case errors.Is(err, errs.ErrInvalidDueDate):
			c.String(http.StatusBadRequest, err.Error())
		default:
			c.String(http.StatusInternalServerError, "internal server error")
	}
}
