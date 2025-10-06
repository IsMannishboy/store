create table cart (
    id serial primary key,
    user_id int not null REFERENCES users(id),
    product_id int not null REFERENCES products(id),
    quantity int default 1,
    created_at timestamp default now()
);
