CREATE TABLE game_saves (
    user_id TEXT NOT NULL,
    save_id TEXT NOT NULL,
    balance BIGINT NOT NULL,
    last_update_at BIGINT NOT NULL,
    miners JSONB,
    PRIMARY KEY (user_id, save_id)
);