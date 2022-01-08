CREATE TABLE IF NOT EXISTS Users (
    id BIGSERIAL PRIMARY KEY,
    fullName TEXT,
    email TEXT UNIQUE,
    pwd TEXT,
    createdAt timestamptz DEFAULT now(),
    updatedAt timestamptz DEFAULT now()
);