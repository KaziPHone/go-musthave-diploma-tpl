CREATE TABLE orders(
    number varchar(50) NOT NULL,
    accrual double precision,
    user_id integer,
    status varchar(20) NOT NULL,
    uploaded_at timestamp without time zone NOT NULL,
    processed_at timestamp without time zone,
    PRIMARY KEY(number),
    CONSTRAINT orders_user_id_fkey FOREIGN key(user_id) REFERENCES users(id)
);