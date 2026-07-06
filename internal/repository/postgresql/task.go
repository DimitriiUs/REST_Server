package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"REST_Server/internal/errors"
	"REST_Server/internal/model"
)

type PgxPoolIFace interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Ping(ctx context.Context) error
	Close()
}
type repo struct {
	db PgxPoolIFace
}

func NewRepo(db PgxPoolIFace) *repo {
	return &repo{db: db}

}

func (r *repo) GetAllTasks() ([]model.Task, error) {
	query := "SELECT * FROM tasks"

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allTasks := make([]model.Task, 0, 100)
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

	row := r.db.QueryRow(context.Background(),
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

	row := r.db.QueryRow(context.Background(),
		query,
		description,
		due.Format(time.DateTime))
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r *repo) DeleteTaskByID(id int) error {
	query := "DELETE FROM tasks WHERE task_id = $1 RETURNING task_id"

	row := r.db.QueryRow(context.Background(),
		query,
		id)
	var ids int
	if err := row.Scan(&ids); err != nil {
		return err
	}
	if ids == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *repo) DeleteAllTasks() error {
	_, err := r.db.Exec(context.Background(), "TRUNCATE TABLE tasks")
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetTaskByDueDate(due time.Time) ([]model.Task, error) {
	query := "SELECT * FROM tasks WHERE due_date = $1"

	rows, err := r.db.Query(context.Background(),
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

	return tasks, nil
}
