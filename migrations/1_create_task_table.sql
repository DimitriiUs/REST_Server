-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_description TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP Table tasks;
-- +goose StatementEnd