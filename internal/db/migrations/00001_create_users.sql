-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    uuid          UUID NOT NULL DEFAULT gen_random_uuid(),
    username      VARCHAR(255) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    created_by    BIGINT REFERENCES users(id),
    updated_by    BIGINT REFERENCES users(id),
    deleted_by    BIGINT REFERENCES users(id)
);

CREATE INDEX idx_users_uuid_hash ON users USING HASH (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
