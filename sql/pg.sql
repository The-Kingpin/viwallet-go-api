CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email character varying(64) NOT NULL UNIQUE,
    password character varying(256) NOT NULL,
    phone_number character varying(64) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
