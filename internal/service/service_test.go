package service

import (
	"REST_Server/internal/errors"
	"REST_Server/internal/model"
	"testing"
	"time"
)

type mockTaskRepository struct {
	testTasks    map[int]model.Task
	nextID       int
	getByIDError error
}

func (m *mockTaskRepository) GetAllTasks() ([]model.Task, error) {
	result := make([]model.Task, 0, len(m.testTasks))

	for _, v := range m.testTasks {
		result = append(result, v)
	}

	return result, nil
}

func (m *mockTaskRepository) GetTaskByID(id int) (model.Task, error) {
	if m.getByIDError != nil {
		return model.Task{}, m.getByIDError
	}

	task, ok := m.testTasks[id]
	if !ok {
		return model.Task{}, errors.ErrNotFound
	}
	return task, nil
}

func (m *mockTaskRepository) CreateTask(description string, due time.Time) (int, error) {
	task := model.Task{
		ID:   m.nextID,
		Text: description,
		Due:  due,
	}

	m.testTasks[m.nextID] = task
	m.nextID++
	return task.ID, nil
}

func (m *mockTaskRepository) DeleteTaskByID(id int) (string, error) {
	task
}

func (m *mockTaskRepository) DeleteAllTasks() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockTaskRepository) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	//TODO implement me
	panic("implement me")
}

func TestGetAllTasks(t *testing.T) {
	taskRepo := &mockTaskRepository{}
	service := NewService(taskRepo)
}
