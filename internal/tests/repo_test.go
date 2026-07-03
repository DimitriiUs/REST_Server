package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/require"

	"REST_Server/internal/repository/postgresql"
)

func TestRepoGetAllTasks_Empty(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	rows := mock.NewRows([]string{"task_id", "task_description", "due_date"})
	mock.ExpectQuery("SELECT (.+) FROM tasks").WillReturnRows(rows)

	testRepo := postgresql.NewRepo(mock)
	tasks, err := testRepo.GetAllTasks()

	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestRepoGetAllTasks_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery("SELECT (.+) FROM tasks").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	_, err = testRepo.GetAllTasks()

	require.Error(t, err, "fail")
}

func TestRepoGetAllTasks_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	rows := mock.NewRows([]string{"task_id", "task_description", "due_date"}).
		AddRow(101, "test task101", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)).
		AddRow(102, "test task102", time.Date(2026, time.June, 17, 0, 0, 0, 0, time.UTC)).
		AddRow(103, "test task103", time.Date(2026, time.June, 17, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT (.+) FROM tasks").WillReturnRows(rows)

	testRepo := postgresql.NewRepo(mock)
	tasks, err := testRepo.GetAllTasks()

	require.NoError(t, err)
	require.Equal(t, 3, len(tasks))
}

func TestRepoGetTaskById_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery("SELECT (.+)").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	_, err = testRepo.GetTaskByID(4)

	require.Error(t, err, "fail")
}

func TestRepoGetTaskById_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	rows := mock.NewRows([]string{"task_id", "task_description", "due_date"}).
		AddRow(101, "test task101", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT (.+)").WithArgs(101).WillReturnRows(rows)

	testRepo := postgresql.NewRepo(mock)
	task, err := testRepo.GetTaskByID(101)

	require.NoError(t, err)
	require.Equal(t, 101, task.ID)
	require.Equal(t, "test task101", task.Text)
	require.Equal(t, time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC), task.Due)
}

func TestRepoCreateTask_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery("INSERT (.+)").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	_, err = testRepo.CreateTask("test task", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC))

	require.Error(t, err, "fail")
}

func TestRepoCreateTask_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	row := mock.NewRows([]string{"task_id"}).
		AddRow(111)
	mock.ExpectQuery("INSERT (.+)").
		WithArgs("test task", "2026-06-15 00:00:00").
		WillReturnRows(row)

	testRepo := postgresql.NewRepo(mock)
	id, err := testRepo.CreateTask("test task", time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC))

	require.NoError(t, err)
	require.Equal(t, 111, id)
}

func TestRepoDeleteTaskByID_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery("DELETE FROM tasks").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	err = testRepo.DeleteTaskByID(4)

	require.Error(t, err, "fail")
}

func TestRepoDeleteTaskByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	row := mock.NewRows([]string{"task_id"}).
		AddRow(99)
	mock.ExpectQuery("DELETE FROM tasks").
		WithArgs(99).
		WillReturnRows(row)

	testRepo := postgresql.NewRepo(mock)
	err = testRepo.DeleteTaskByID(99)

	require.NoError(t, err)
}

func TestRepoDeleteAllTasks_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectExec("TRUNCATE TABLE tasks").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	err = testRepo.DeleteAllTasks()

	require.Error(t, err, "fail")
}

func TestRepoGetTaskByDueDate_Fail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	mock.ExpectQuery("SELECT (.+)").WillReturnError(errors.New("fail"))

	testRepo := postgresql.NewRepo(mock)
	dueDate := time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)
	_, err = testRepo.GetTaskByDueDate(dueDate)

	require.Error(t, err, "fail")
}

func TestRepoGetTaskByDueDate_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	date := time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)
	rows := mock.NewRows([]string{"task_id", "task_description", "due_date"}).
		AddRow(111, "test task111", date).
		AddRow(102, "test task102", date).
		AddRow(112, "test task112", date)
	mock.ExpectQuery("SELECT (.+) FROM tasks").
		WithArgs(date).
		WillReturnRows(rows)

	testRepo := postgresql.NewRepo(mock)
	tasks, err := testRepo.GetTaskByDueDate(date)

	require.NoError(t, err)
	require.Equal(t, 3, len(tasks))
}
