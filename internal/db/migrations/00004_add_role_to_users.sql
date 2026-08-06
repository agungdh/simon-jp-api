-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'pegawai'
        CONSTRAINT chk_users_role CHECK (role IN ('admin', 'pegawai'));

UPDATE users SET role = 'admin' WHERE username = 'admin';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS role;
-- +goose StatementEnd
