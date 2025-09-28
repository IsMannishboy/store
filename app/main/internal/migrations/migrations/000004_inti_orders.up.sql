create table orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    total_amount DECIMAL(10,2),
    status VARCHAR(50) DEFAULT 'pending', -- pending, paid, shipped, cancelled
    created_at TIMESTAMP DEFAULT NOW()
)
