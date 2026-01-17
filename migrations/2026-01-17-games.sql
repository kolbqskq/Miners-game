CREATE TABLE games (
    user_id TEXT NOT NULL,
    game_id TEXT NOT NULL,
    balance BIGINT NOT NULL,
    income BIGINT NOT NULL,
    last_update_at BIGINT NOT NULL,
    miners JSONB,
    equipments JSONB,
    upgrades JSONB,
    PRIMARY KEY (user_id, game_id)
);