CREATE TABLE users (
    user_id TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    PRIMARY KEY (user_id),
    UNIQUE (email, username)
);