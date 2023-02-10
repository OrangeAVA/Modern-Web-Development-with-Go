CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    username text NOT NULL UNIQUE,
    user_password text NOT NULL,
    user_role text NOT NULL,
    access_token text,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

CREATE INDEX user_access_token
ON users (access_token);

INSERT INTO users(username, user_password, user_role)
VALUES
    ('admin', crypt('admin', gen_salt('bf')), 'admin'),
    ('runner', crypt('runner', gen_salt('bf')), 'runner');
