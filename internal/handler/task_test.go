package handler_test

import (
	"encoding/json"
	errors2 "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"REST_Server/internal/errors"
	"REST_Server/internal/handler"
	"REST_Server/internal/handler/fake"
	"REST_Server/internal/model"
)

var (
	fakeTestService = &fake.FakeTaskService{}
	testHandler     = handler.NewHandler(fakeTestService)
	router          = setupRouter()
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	r := gin.New()
	r.GET("/tasks", testHandler.GetAllTasks)
	r.POST("/task", testHandler.CreateTask)
	r.GET("/task/:id", testHandler.GetTaskByID)
	r.DELETE("/task/:id", testHandler.DeleteTaskByID)
	r.DELETE("/tasks", testHandler.DeleteAllTasks)
	r.GET("/task/due/:year/:month/:day", testHandler.GetTasksByDue)
	return r
}

func testRequestResponse(path string, method string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func TestHandlerGetAllTasks_Empty(t *testing.T) {
	result := errors.ErrNotFound
	fakeTestService.GetAllTasksReturns(nil, result)

	resp := testRequestResponse("/tasks", "GET", nil)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandlerGetAllTasks_Success(t *testing.T) {
	resultFakeService := []model.Task{
		{ID: 1, Text: "test1", Due: time.Date(2026, time.April, 2, 15, 4, 5, 0, time.UTC)},
		{ID: 2, Text: "test2", Due: time.Date(2026, time.April, 3, 17, 4, 5, 0, time.UTC)},
		{ID: 3, Text: "test3", Due: time.Date(2026, time.April, 4, 18, 4, 5, 0, time.UTC)},
	}
	fakeTestService.GetAllTasksReturns(resultFakeService, nil)

	resp := testRequestResponse("/tasks", "GET", nil)

	var responseResult []model.Task
	err := json.Unmarshal(resp.Body.Bytes(), &responseResult)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, len(resultFakeService), len(responseResult))
}

func TestHandlerGetTaskByID_NotFound(t *testing.T) {
	result1 := model.Task{}
	result2 := errors.ErrNotFound
	fakeTestService.GetTaskByIDReturns(result1, result2)

	resp := testRequestResponse("/task/99", "GET", nil)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandlerGetTaskByID_InvalidID(t *testing.T) {
	result1 := model.Task{}
	result2 := errors.ErrInvalidID
	fakeTestService.GetTaskByIDReturns(result1, result2)

	resp := testRequestResponse("/task/т", "GET", nil)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHandlerGetTaskByID_Success(t *testing.T) {
	result := model.Task{
		ID:   112,
		Text: "test1",
		Due:  time.Date(2026, time.April, 2, 15, 4, 5, 0, time.UTC),
	}
	fakeTestService.GetTaskByIDReturns(result, nil)

	resp := testRequestResponse("/task/112", "GET", nil)

	var responseResult model.Task
	err := json.Unmarshal(resp.Body.Bytes(), &responseResult)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, result, responseResult)
}

func TestHandlerGetTaskByID_Fail(t *testing.T) {
	fakeTestService.GetTaskByIDReturns(model.Task{}, pgx.ErrNoRows)

	resp := testRequestResponse("/task/121", "GET", nil)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestHandlerCreateTask_InvalidDescription(t *testing.T) {
	fakeTestService.CreateTaskReturns(0, errors.ErrInvalidDescription)

	task := `{"text": "",	"due": "2026-05-19T00:00:00Z"}`
	resp := testRequestResponse("/task", "POST", strings.NewReader(task))

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHandlerCreateTask_InvalidDue(t *testing.T) {
	fakeTestService.CreateTaskReturns(0, errors.ErrInvalidDueDate)

	task := `{"text": "test task",	"due": ""}`
	resp := testRequestResponse("/task", "POST", strings.NewReader(task))

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHandlerCreateTask_Success(t *testing.T) {
	fakeTestService.CreateTaskReturns(101, nil)

	task := `{"text": "test task",	"due": "2026-05-19T00:00:00Z"}`
	resp := testRequestResponse("/task", "POST", strings.NewReader(task))

	var responseTask model.Task

	err := json.Unmarshal(resp.Body.Bytes(), &responseTask)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)
}

func TestHandlerDeleteTaskByID_InvalidID(t *testing.T) {
	fakeTestService.DeleteTaskByIDReturns(errors.ErrInvalidID)

	resp := testRequestResponse("/task/99", "DELETE", nil)

	require.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestHandlerDeleteTaskByID_NotFound(t *testing.T) {
	fakeTestService.DeleteTaskByIDReturns(errors.ErrNotFound)

	resp := testRequestResponse("/task/100", "DELETE", nil)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandlerDeleteTaskByID_Success(t *testing.T) {
	fakeTestService.DeleteTaskByIDReturns(nil)

	resp := testRequestResponse("/task/101", "DELETE", nil)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func TestHandlerDeleteAllTasks_Success(t *testing.T) {
	fakeTestService.DeleteAllTasksReturns(nil)

	resp := testRequestResponse("/tasks", "DELETE", nil)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func TestHandlerDeleteAllTasks_Fail(t *testing.T) {
	fakeTestService.DeleteAllTasksReturns(errors2.New("test error"))

	resp := testRequestResponse("/tasks", "DELETE", nil)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestHandlerGetTasksByDue_InvalidDue(t *testing.T) {
	fakeTestService.GetTasksByDueReturns(nil, errors.ErrInvalidDueDate)

	resp1 := testRequestResponse("/task/due/г/5/25", "GET", nil)
	resp2 := testRequestResponse("/task/due/2026/м/25", "GET", nil)
	resp3 := testRequestResponse("/task/due/2026/5/д", "GET", nil)

	require.Equal(t, http.StatusBadRequest, resp1.Code)
	require.Equal(t, http.StatusBadRequest, resp2.Code)
	require.Equal(t, http.StatusBadRequest, resp3.Code)
}

func TestHandlerGetTasksByDue_Success(t *testing.T) {
	result := []model.Task{
		{ID: 112, Text: "test112", Due: time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC)},
		{ID: 222, Text: "test222", Due: time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC)},
		{ID: 335, Text: "test335", Due: time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC)},
	}
	fakeTestService.GetTasksByDueReturns(result, nil)

	resp := testRequestResponse("/task/due/2026/6/10", "GET", nil)

	var responseResult []model.Task
	err := json.Unmarshal(resp.Body.Bytes(), &responseResult)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, result, responseResult)
}
