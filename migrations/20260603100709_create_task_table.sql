-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
    task_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_description TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL
);

-- +goose Down
DROP Table tasks;
