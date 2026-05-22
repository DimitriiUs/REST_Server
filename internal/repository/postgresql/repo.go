package postgresql

import (
	"REST_Server/internal/model"
	"time"
)

type repo struct {
}

func (repo) GetAllTasks() ([]model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (repo) GetTaskByID(id int) (model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (repo) CreateTask(description string, due time.Time) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (repo) DeleteTaskByID(id int) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (repo) DeleteAllTasks() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (repo) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	//TODO implement me
	panic("implement me")
}
