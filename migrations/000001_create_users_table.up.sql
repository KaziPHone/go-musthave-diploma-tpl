CREATE TABLE users(
    id integer GENERATED ALWAYS AS IDENTITY NOT NULL,
    create_time timestamp without time zone DEFAULT '2026-01-10 19:28:13.459931'::timestamp without time zone,
    login varchar(255),
    password varchar(255),
    PRIMARY KEY(id)
);