SET TIMEZONE="Europe/Moscow";

CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    login VARCHAR(255) NOT NULL,
    salt BYTEA NOT NULL,
    password TEXT NOT NULL,
    name VARCHAR(255),
    surname VARCHAR(255),
    birthdate DATE,
    email VARCHAR(255),
    phone VARCHAR(32),
    PRIMARY KEY (id),
    UNIQUE (login)
);
