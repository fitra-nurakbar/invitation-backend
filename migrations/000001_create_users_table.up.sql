CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    password   TEXT NOT NULL,
    role       user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);