CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    username   VARCHAR(50) UNIQUE NOT NULL,
    password   VARCHAR(64) NOT NULL,
    salt       VARCHAR(64)  NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
