CREATE TABLE
    IF NOT EXISTS users (
        id UUID PRIMARY KEY NOT NULL,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        hashed_password TEXT NOT NULL,
        fullname TEXT,
        status TEXT,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP,
        deleted_at TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS jokes (
        id UUID PRIMARY KEY NOT NULL,
        title TEXT NOT NULL,
        text TEXT NOT NULL,
        explanation TEXT,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP,
        deleted_at TIMESTAMP,
        user_id UUID NOT NULL,
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_joke_title ON jokes (user_id, title);