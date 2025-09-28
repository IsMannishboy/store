create table products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150),
    description TEXT,
    price DECIMAL(10,2),
    stock INT,
    category_id INT REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW()
)
