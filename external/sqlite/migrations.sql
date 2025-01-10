CREATE TABLE
    users (
        id TEXT PRIMARY KEY,
        username TEXT NOT NULL,
        email TEXT NOT NULL,
        hashed_password TEXT NOT NULL,
        fullname TEXT,
        status TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL,
        deleted_at DATETIME
    );

CREATE TABLE
    jokes (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        text TEXT NOT NULL,
        explanation TEXT,
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL,
        deleted_at DATETIME,
        user_id TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );