create table products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150),
    description TEXT,
    price DECIMAL(10,2),
    stock INT,
    category VARCHAR(100) REFERENCES categories(c_name),
    created_at TIMESTAMP DEFAULT NOW()
)
