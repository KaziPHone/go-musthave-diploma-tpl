CREATE TABLE balance_operations(
    id SERIAL NOT NULL,
    user_id integer,
    "type" varchar(20) NOT NULL,
    amount double precision NOT NULL,
    order_number varchar(50),
    processed_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    CONSTRAINT balance_operations_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);