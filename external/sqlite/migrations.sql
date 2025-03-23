CREATE TABLE
    users (
        id TEXT PRIMARY KEY NOT NULL,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        hashed_password TEXT NOT NULL,
        fullname TEXT,
        status TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME,
        deleted_at DATETIME
    );

CREATE TABLE
    jokes (
        id TEXT PRIMARY KEY NOT NULL,
        title TEXT NOT NULL,
        text TEXT NOT NULL,
        explanation TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME,
        deleted_at DATETIME,
        user_id TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE UNIQUE INDEX idx_user_joke_title ON jokes (user_id, title);