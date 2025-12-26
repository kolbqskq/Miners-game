CREATE TABLE sessions (
    id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    expires_at BIGINT NOT NULL,
    PRIMARY KEY (user_id, id)
);