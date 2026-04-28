package taskstore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Due  time.Time `json:"due"`
}

func CreateTask(text string, due time.Time) int {

	row := pool.QueryRow(context.Background(), "INSERT INTO tasks (task_description, due_date) VALUES ($1, $2) RETURNING task_id", text, due.Format(time.DateTime))
	var count int
	if err := row.Scan(&count); err != nil {
		log.Fatal(err)

	}

	return count
}

func GetTask(id int) (Task, error) {
	task := Task{}
	row := pool.QueryRow(context.Background(), "SELECT * FROM tasks WHERE task_id = $1", id)
	if err := row.Scan(&task.Id, &task.Text, &task.Due); err != nil {
		return Task{}, err
	}

	return task, nil
}

func DeleteTask(id int) (string, error) {
	row := pool.QueryRow(context.Background(), "DELETE FROM tasks WHERE task_id = $1 RETURNING task_description", id)
	var description string
	if err := row.Scan(&description); err != nil {
		return "", err
	}

	return fmt.Sprintf("Task: `%s` was deleted", description), nil
}

func DeleteAllTasks() (string, error) {
	res, err := pool.Exec(context.Background(), "TRUNCATE TABLE tasks")
	if err != nil {
		return "", err
	}

	if res.RowsAffected() == 0 {
		return "", errors.New("No tasks were deleted")
	}
	return fmt.Sprintf("Deleted %d tasks", res.RowsAffected()), nil
}

func GetAllTasks() ([]Task, error) {
	rows, err := pool.Query(context.Background(), "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allTasks []Task
	for rows.Next() {
		task := Task{}
		if err := rows.Scan(&task.Id, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		allTasks = append(allTasks, task)
	}
	return allTasks, nil
}

func GetTasksByDueDate(year int, month time.Month, day int) ([]Task, error) {
	dueDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	rows, err := pool.Query(context.Background(), "SELECT * FROM tasks WHERE due_date = $1", dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}
		if err := rows.Scan(&task.Id, &task.Text, &task.Due); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if tasks == nil {
		return nil, errors.New("No tasks were found")
	}
	return tasks, nil
}
