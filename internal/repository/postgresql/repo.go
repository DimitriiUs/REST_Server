package postgresql

import (
	"REST_Server/internal/model"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *repo {
	return &repo{pool: pool}
}

func (r *repo) GetAllTasks() ([]model.Task, error) {
	rows, err := r.pool.Query(context.Background(), "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []model.Task
	for rows.Next() {
		task := model.Task{}
		if err := rows.Scan(&task.Id, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

func (r *repo) GetTaskByID(id int) (model.Task, error) {
	task := model.Task{}
	row := r.pool.QueryRow(context.Background(), "SELECT * FROM tasks WHERE task_id = $1", id)
	if err := row.Scan(&task.Id, &task.Text, &task.Due); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *repo) CreateTask(description string, due time.Time) (int, error) {
	row := r.pool.QueryRow(context.Background(),
		"INSERT INTO tasks (task_description, due_date) VALUES ($1, $2) RETURNING task_id",
		description,
		due.Format(time.DateTime))
	var id int
	if err := row.Scan(&id); err != nil {
		log.Fatal(err)
		return -1, err
	}

	return id, nil
}

func (r *repo) DeleteTaskByID(id int) (string, error) {
	row := r.pool.QueryRow(context.Background(),
		"DELETE FROM tasks WHERE task_id = $1 RETURNING task_description",
		id)
	var description string
	if err := row.Scan(&description); err != nil {
		return "", err
	}

	return fmt.Sprintf("Task: `%s` was deleted", description), nil
}

func (r *repo) DeleteAllTasks() (string, error) {
	res, err := r.pool.Exec(context.Background(), "TRUNCATE TABLE tasks")
	if err != nil {
		return "", err
	}

	if res.RowsAffected() == 0 {
		return "", errors.New("no tasks were deleted")
	}
	return fmt.Sprintf("Deleted %d tasks", res.RowsAffected()), nil
}

func (r *repo) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT * FROM tasks WHERE due_date = $1",
		due)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		task := model.Task{}
		if err := rows.Scan(&task.Id, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		return nil, errors.New("no tasks were found")
	}
	return tasks, nil
}
