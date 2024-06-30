BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    passportNumber VARCHAR(255) NOT NULL,
    surname VARCHAR(255),
    name VARCHAR(255),
    patronymic VARCHAR(255),
    address TEXT
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    userId INT REFERENCES users(id) ON DELETE CASCADE,
    taskName VARCHAR(255) NOT NULL,
    content TEXT,
    startTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    endTime TIMESTAMP
);

COMMIT;