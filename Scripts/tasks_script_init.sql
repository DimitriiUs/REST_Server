CREATE TABLE IF NOT EXISTS tasks (
    task_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_description TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL
);

--TRUNCATE TABLE tasks;

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