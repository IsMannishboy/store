create table users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);