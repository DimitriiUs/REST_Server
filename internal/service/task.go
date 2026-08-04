package service

import (
	"log/slog"
	"strconv"
	"time"

	errs "REST_Server/internal/errors"
	"REST_Server/internal/model"
)

//go:generate go tool counterfeiter -o /fake . TaskRepository

type TaskRepository interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskByID(id int) (*model.Task, error)
	CreateTask(description string, due time.Time) (int, error)
	DeleteTaskByID(id int) error
	DeleteAllTasks() error
	GetTaskByDueDate(due time.Time) ([]model.Task, error)
}

type service struct {
	repo TaskRepository
	log  *slog.Logger
}

func NewService(repo TaskRepository, log *slog.Logger) *service {
	return &service{repo,log}
}

func (s *service) GetAllTasks() ([]model.Task, error) {
	tasks, err := s.repo.GetAllTasks()
	if len(tasks) == 0 {
		s.log.Info(errs.ErrNotFound.Error())
		return nil, errs.ErrNotFound
	}
	return tasks, err
}

func (s *service) GetTaskByID(ids string) (model.Task, error) {
	id, err := strconv.Atoi(ids)
	if err != nil {
		s.log.Error("%w: %s - not a number",errs.ErrInvalidID.Error(), ids)
		return model.Task{}, errs.ErrInvalidID
	}
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		s.log.Error(err.Error())
		return model.Task{}, err
	}

	if task.IsEmpty() {
		s.log.Error("%w: %d", errs.ErrNotFound.Error(), id)
		return *task, errs.ErrNotFound
	}
	return *task, err
}

func (s *service) CreateTask(description string, due time.Time) (int, error) {
	if description == "" {
		s.log.Error("%w: empty description",errs.ErrInvalidDescription.Error(), description)
		return 0, errs.ErrInvalidDescription
	}
	if due.IsZero() {
		s.log.Error("%w: empty due date", errs.ErrInvalidDueDate.Error(), due)
		return 0, errs.ErrInvalidDueDate
	}
	return s.repo.CreateTask(description, due)
}

func (s *service) DeleteTaskByID(ids string) error {
	id, err := strconv.Atoi(ids)
	if err != nil || id == 0 {
		s.log.Error("%w: %s - not a number",errs.ErrInvalidID.Error(), ids)
		return errs.ErrInvalidID
	}
	return s.repo.DeleteTaskByID(id)
}

func (s *service) DeleteAllTasks() error {
	return s.repo.DeleteAllTasks()
}

func (s *service) GetTasksByDue(year string, month string, day string) ([]model.Task, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		s.log.Error("%w: %s - not valid year", errs.ErrInvalidDueDate.Error(), year)
		return nil, errs.ErrInvalidDueDate
	}

	intMonth, err := strconv.Atoi(month)
	if err != nil || intMonth < int(time.January) || intMonth > int(time.December) {
		s.log.Error("%w: %s - not valid month", errs.ErrInvalidDueDate.Error(), month)
		return nil, errs.ErrInvalidDueDate
	}

	intDay, err := strconv.Atoi(day)
	if err != nil {
		return nil, errs.ErrInvalidDueDate
	}
	dueDate := time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)

	tasks, err := s.repo.GetTaskByDueDate(dueDate)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	if len(tasks) == 0 {
		s.log.Info(errs.ErrNotFound.Error())
		return nil, errs.ErrNotFound
	}

	return tasks, nil
}
