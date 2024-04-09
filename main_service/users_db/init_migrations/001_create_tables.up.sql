CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SET TIMEZONE="Europe/Moscow";

CREATE TABLE credentials (
    id uuid NOT NULL,
    login character varying(255) NOT NULL,
    salt bytea NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (login)
);

CREATE TABLE personal_data (
    id uuid NOT NULL,
    name character varying(255),
    surname character varying(255),
    birthdate date,
    email character varying(255),
    phone character varying(32),
    PRIMARY KEY (id)
);
