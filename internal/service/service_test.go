package service

import (
	"REST_Server/internal/errors"
	"REST_Server/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func (m *mockTaskRepository) DeleteTaskByID(id int) error {
	_, err := m.GetTaskByID(id)
	if err != nil {
		return err
	}
	delete(m.testTasks, id)
	return nil
}

func (m *mockTaskRepository) DeleteAllTasks() error {
	clear(m.testTasks)
	return nil
}

func (m *mockTaskRepository) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	var result []model.Task

	for _, v := range m.testTasks {
		if v.Due == due {
			result = append(result, v)
		}
	}

	return result, nil
}

func TestGetAllTasks_Empty(t *testing.T) {
	taskRepo := &mockTaskRepository{
		testTasks: make(map[int]model.Task),
		nextID:    1,
	}
	service := NewService(taskRepo)

	tasks, err := service.GetAllTasks()

	require.Empty(t, tasks)
	require.NoError(t, err)

}

func TestGetAllTasks_Success(t *testing.T) {
	taskRepo := &mockTaskRepository{
		testTasks: map[int]model.Task{
			1: model.Task{
				ID:   1,
				Text: "test",
				Due:  time.Now(),
			},
			2: model.Task{
				ID:   2,
				Text: "test1",
				Due:  time.Now(),
			},
			3: model.Task{
				ID:   3,
				Text: "test2",
				Due:  time.Now(),
			},
		},
		nextID: 3,
	}

	service := NewService(taskRepo)
	tasks, err := service.GetAllTasks()

	require.NoError(t, err)
	require.Len(t, tasks, 3)
}

func TestGetTaskByID_InvalidID(t *testing.T) {
	taskRepo := &mockTaskRepository{
		testTasks: map[int]model.Task{
			1: model.Task{ID: 1, Text: "test", Due: time.Now()},
			2: model.Task{ID: 2, Text: "test1", Due: time.Now()},
			3: model.Task{ID: 3, Text: "test2", Due: time.Now()},
		},
		nextID: 3,
	}

	service := NewService(taskRepo)

	_, err := service.GetTaskByID("т")

	require.ErrorIs(t, err, errors.ErrInvalidID)
}

func TestGetTaskByID_NotFound(t *testing.T) {

}

func TestCreateTask_InvalidID(t *testing.T) {
	taskRepo := &mockTaskRepository{
		testTasks: make(map[int]model.Task),
		nextID:    7,
	}
	service := NewService(taskRepo)

	id, err := service.CreateTask("test task", time.Now())

	require.Error(t, err)

}
