package tests

import (
	"REST_Server/internal/errors"
	"REST_Server/internal/model"
	"REST_Server/internal/service"
	"REST_Server/internal/tests/fakes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	testRepo    = &fakes.FakeTaskRepository{}
	testService = service.NewService(testRepo)
)

func TestGetAllTasks_Empty(t *testing.T) {
	testRepo.GetAllTasksReturns([]model.Task{}, nil)

	tasks, err := testService.GetAllTasks()

	require.Nil(t, tasks)
	require.Error(t, err, errors.ErrNotFound)

}

func TestGetAllTasks_Success(t *testing.T) {
	result := []model.Task{
		{ID: 1, Text: "test1", Due: time.Now()},
		{ID: 2, Text: "test2", Due: time.Now()},
		{ID: 3, Text: "test3", Due: time.Now()},
	}
	testRepo.GetAllTasksReturns(result, nil)

	tasks, err := testService.GetAllTasks()

	require.NoError(t, err)
	require.Len(t, tasks, 3)
}

func TestGetTaskByID_InvalidID(t *testing.T) {
	_, err := testService.GetTaskByID("т")

	require.ErrorIs(t, err, errors.ErrInvalidID)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	testRepo.GetTaskByIDReturns(model.Task{}, nil)

	_, err := testService.GetTaskByID("5")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestCreateTask_InvalidDescription(t *testing.T) {
	_, err := testService.CreateTask("", time.Now())

	require.Error(t, err, errors.ErrInvalidDescription)
}

func TestCreateTask_InvalidDueDate(t *testing.T) {
	_, err := testService.CreateTask("test task", time.Time{})

	require.Error(t, err, errors.ErrInvalidDueDate)
}

func TestCreateTask_Success(t *testing.T) {
	testRepo.CreateTaskReturns(4, nil)

	id, err := testService.CreateTask("test task", time.Now())

	require.NoError(t, err)
	require.Equal(t, id, 4)
}

func TestDeleteTaskByID_InvalidID(t *testing.T) {
	err := testService.DeleteTaskByID("т")

	require.Error(t, err, errors.ErrInvalidID)
}

func TestDeleteTaskByID_NotFound(t *testing.T) {
	testRepo.DeleteTaskByIDReturns(errors.ErrNotFound)
	err := testService.DeleteTaskByID("5")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestDeleteTaskByID_Success(t *testing.T) {
	testRepo.DeleteTaskByIDReturns(nil)
	err := testService.DeleteTaskByID("1")

	require.NoError(t, err)
}

func TestDeleteAllTasks_Success(t *testing.T) {
	testRepo.DeleteAllTasksReturns(nil)
	err := testService.DeleteAllTasks()

	require.NoError(t, err)
}

func TestGetTaskByDueDate_InvalidDueDate(t *testing.T) {
	testRepo.GetTaskByDueDateReturns(nil, errors.ErrInvalidDueDate)

	_, errYear := testService.GetTasksByDue("year", "05", "20")
	_, errMonth := testService.GetTasksByDue("2026", "13", "02")
	_, errDay := testService.GetTasksByDue("2026", "05", "day")
	require.Error(t, errYear, errors.ErrInvalidDueDate)
	require.Error(t, errMonth, errors.ErrInvalidDueDate)
	require.Error(t, errDay, errors.ErrInvalidDueDate)
}

func TestGetTaskByDueDay_NotFound(t *testing.T) {
	testRepo.GetTaskByDueDateReturns(nil, errors.ErrNotFound)

	_, err := testService.GetTasksByDue("2026", "05", "20")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestGetTaskByDueDay_Success(t *testing.T) {
	result := []model.Task{
		{ID: 2, Text: "test2", Due: time.Date(2026, time.May, 22, 0, 0, 0, 0, time.UTC)},
		{ID: 3, Text: "test3", Due: time.Date(2026, time.May, 22, 0, 0, 0, 0, time.UTC)},
	}
	testRepo.GetTaskByDueDateReturns(result, nil)
	tasks, err := testService.GetTasksByDue("2026", "05", "22")

	require.NoError(t, err)
	require.Len(t, tasks, 2)
}
