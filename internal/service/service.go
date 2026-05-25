package service

import (
	"REST_Server/internal/model"
	"errors"
	"log"
	"strconv"
	"time"
)

type TaskRepository interface {
	GetAllTasks() ([]model.Task, error)
	GetTaskByID(id int) (model.Task, error)
	CreateTask(description string, due time.Time) (int, error)
	DeleteTaskByID(id int) (string, error)
	DeleteAllTasks() (string, error)
	GetTaskByDueDate(due time.Time) ([]model.Task, error)
}

type service struct {
	repo TaskRepository
}

func NewService(repo TaskRepository) *service {
	return &service{repo}
}

func (s *service) GetAllTasks() ([]model.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *service) GetTaskByID(ids string) (model.Task, error) {
	id, err := strconv.Atoi(ids)
	if err != nil {
		log.Println(err)
		return model.Task{}, err
	}
	return s.repo.GetTaskByID(id)
}

func (s *service) CreateTask(description string, due time.Time) (int, error) {
	if description == "" {
		//log
		return 0, errors.New("description is required")
	}
	if due.IsZero() {
		//log
		return 0, errors.New("due is required")
	}
	return s.repo.CreateTask(description, due)
}

func (s *service) DeleteTaskByID(ids string) (string, error) {
	id, err := strconv.Atoi(ids)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return s.repo.DeleteTaskByID(id)
}

func (s *service) DeleteAllTasks() (string, error) {
	return s.repo.DeleteAllTasks()
}

func (s *service) GetTasksByDue(year string, month string, day string) ([]model.Task, error) {
	intYear, err := strconv.Atoi(year)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	intMonth, err := strconv.Atoi(month)
	if err != nil || intMonth < int(time.January) || intMonth > int(time.December) {
		log.Println(err)
		return nil, err
	}

	intDay, err := strconv.Atoi(day)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dueDate := time.Date(intYear, time.Month(intMonth), intDay, 0, 0, 0, 0, time.UTC)

	return s.repo.GetTaskByDueDate(dueDate)
}
