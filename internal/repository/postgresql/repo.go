package postgresql

import (
	"REST_Server/internal/errors"
	"REST_Server/internal/model"

	"context"
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
	query := "SELECT * FROM tasks"

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []model.Task
	for rows.Next() {
		task := model.Task{}
		if err := rows.Scan(&task.ID, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

func (r *repo) GetTaskByID(id int) (model.Task, error) {
	query := "SELECT * FROM tasks WHERE task_id = $1"
	task := model.Task{}

	row := r.pool.QueryRow(context.Background(),
		query,
		id)
	if err := row.Scan(&task.ID, &task.Text, &task.Due); err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *repo) CreateTask(description string, due time.Time) (int, error) {
	query := `
	INSERT INTO tasks 
    (task_description, due_date) 
	VALUES ($1, $2) 
	RETURNING task_id`

	row := r.pool.QueryRow(context.Background(),
		query,
		description,
		due.Format(time.DateTime))
	var id int
	if err := row.Scan(&id); err != nil {
		log.Fatal(err)
		return -1, err
	}

	return id, nil
}

func (r *repo) DeleteTaskByID(id int) error {
	_, err := r.GetTaskByID(id)
	if err != nil {
		return err
	}

	query := "DELETE FROM tasks WHERE task_id = $1 RETURNING task_description"

	_, err = r.pool.Exec(context.Background(),
		query,
		id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteAllTasks() error {
	_, err := r.pool.Exec(context.Background(), "TRUNCATE TABLE tasks")
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	query := "SELECT * FROM tasks WHERE due_date = $1"

	rows, err := r.pool.Query(context.Background(),
		query,
		due)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		task := model.Task{}
		if err := rows.Scan(&task.ID, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		return nil, errors.ErrNotFound
	}
	return tasks, nil
}
