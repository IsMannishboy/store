create table products (
    id SERIAL PRIMARY KEY,
    prod_name VARCHAR(150),
    description TEXT,
    price DECIMAL(10,2),
    stock INT,
    category VARCHAR(100) ,
    created_at TIMESTAMP DEFAULT NOW()
)
