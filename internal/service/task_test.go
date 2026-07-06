package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"REST_Server/internal/errors"
	"REST_Server/internal/model"
	"REST_Server/internal/service"
	"REST_Server/internal/service/fake"
)

var (
	fakeTestRepo = &fake.FakeTaskRepository{}
	testService  = service.NewService(fakeTestRepo)
)

func TestServiceGetAllTasks_Empty(t *testing.T) {
	fakeTestRepo.GetAllTasksReturns([]model.Task{}, nil)

	tasks, err := testService.GetAllTasks()

	require.Nil(t, tasks)
	require.Error(t, err, errors.ErrNotFound)

}

func TestServiceGetAllTasks_Success(t *testing.T) {
	result := []model.Task{
		{ID: 1, Text: "test1", Due: time.Now()},
		{ID: 2, Text: "test2", Due: time.Now()},
		{ID: 3, Text: "test3", Due: time.Now()},
	}
	fakeTestRepo.GetAllTasksReturns(result, nil)

	tasks, err := testService.GetAllTasks()

	require.NoError(t, err)
	require.Len(t, tasks, 3)
}

func TestServiceGetTaskByID_InvalidID(t *testing.T) {
	_, err := testService.GetTaskByID("т")

	require.ErrorIs(t, err, errors.ErrInvalidID)
}

func TestServiceGetTaskByID_NotFound(t *testing.T) {
	fakeTestRepo.GetTaskByIDReturns(&model.Task{}, nil)

	_, err := testService.GetTaskByID("5")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestServiceCreateTask_InvalidDescription(t *testing.T) {
	_, err := testService.CreateTask("", time.Now())

	require.Error(t, err, errors.ErrInvalidDescription)
}

func TestServiceCreateTask_InvalidDueDate(t *testing.T) {
	_, err := testService.CreateTask("test task", time.Time{})

	require.Error(t, err, errors.ErrInvalidDueDate)
}

func TestServiceCreateTask_Success(t *testing.T) {
	fakeTestRepo.CreateTaskReturns(4, nil)

	id, err := testService.CreateTask("test task", time.Now())

	require.NoError(t, err)
	require.Equal(t, id, 4)
}

func TestServiceDeleteTaskByID_InvalidID(t *testing.T) {
	err := testService.DeleteTaskByID("т")

	require.Error(t, err, errors.ErrInvalidID)
}

func TestServiceDeleteTaskByID_NotFound(t *testing.T) {
	fakeTestRepo.DeleteTaskByIDReturns(errors.ErrNotFound)
	err := testService.DeleteTaskByID("5")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestServiceDeleteTaskByID_Success(t *testing.T) {
	fakeTestRepo.DeleteTaskByIDReturns(nil)
	err := testService.DeleteTaskByID("1")

	require.NoError(t, err)
}

func TestServiceDeleteAllTasks_Success(t *testing.T) {
	fakeTestRepo.DeleteAllTasksReturns(nil)
	err := testService.DeleteAllTasks()

	require.NoError(t, err)
}

func TestServiceGetTaskByDueDate_InvalidDueDate(t *testing.T) {
	fakeTestRepo.GetTaskByDueDateReturns(nil, errors.ErrInvalidDueDate)

	_, errYear := testService.GetTasksByDue("year", "05", "20")
	_, errMonth := testService.GetTasksByDue("2026", "13", "02")
	_, errDay := testService.GetTasksByDue("2026", "05", "day")
	require.Error(t, errYear, errors.ErrInvalidDueDate)
	require.Error(t, errMonth, errors.ErrInvalidDueDate)
	require.Error(t, errDay, errors.ErrInvalidDueDate)
}

func TestServiceGetTaskByDueDay_NotFound(t *testing.T) {
	fakeTestRepo.GetTaskByDueDateReturns(nil, errors.ErrNotFound)

	_, err := testService.GetTasksByDue("2026", "05", "20")

	require.ErrorIs(t, err, errors.ErrNotFound)
}

func TestServiceGetTaskByDueDay_Success(t *testing.T) {
	result := []model.Task{
		{ID: 2, Text: "test2", Due: time.Date(2026, time.May, 22, 0, 0, 0, 0, time.UTC)},
		{ID: 3, Text: "test3", Due: time.Date(2026, time.May, 22, 0, 0, 0, 0, time.UTC)},
	}
	fakeTestRepo.GetTaskByDueDateReturns(result, nil)
	tasks, err := testService.GetTasksByDue("2026", "05", "22")

	require.NoError(t, err)
	require.Len(t, tasks, 2)
}
