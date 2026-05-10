-- +goose Up
-- +goose StatementBegin
DO $$
    BEGIN
        FOR i IN 1..10 LOOP
                INSERT INTO tasks (task_description, due_date)
                VALUES (
                           'Задача №' || i || ': ' || md5(random()::text),
                           NOW() + (random() * (interval '30 days'))
                       );
            END LOOP;
    END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE tasks;
-- +goose StatementEnd