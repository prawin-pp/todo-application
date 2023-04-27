BEGIN;

INSERT INTO
    users (username, password)
VALUES
    (
        'tester01',
        '$2a$10$d0cE9CGFnD0V7lja5zJrN.nBIZiz8cuxjnM6Jrb.tgefvJ2z7BPbW'
    ) ON CONFLICT (username) DO
UPDATE
SET
    password = EXCLUDED.password;

INSERT INTO
    users (username, password)
VALUES
    (
        'tester02',
        '$2a$10$KbV7RmxqCM9vR0id6l6pUu5nkmNt/9tda1ZwqGLrmd7UsYBZNEij2'
    ) ON CONFLICT (username) DO
UPDATE
SET
    password = EXCLUDED.password;

INSERT INTO
    users (username, password)
VALUES
    (
        'tester03',
        '$2a$10$KbV7RmxqCM9vR0id6l6pUu5nkmNt/9tda1ZwqGLrmd7UsYBZNEij2'
    ) ON CONFLICT (username) DO
UPDATE
SET
    password = EXCLUDED.password;

COMMIT;
