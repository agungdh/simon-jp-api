-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, password)
VALUES ('admin', '$2a$10$SFB388dhzc.fHMrJcJ9OKeUkET5ZRXDpOfNGp4cT5dMgluRKnOCay') -- bcrypt("admin123")
ON CONFLICT (username) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE username = 'admin';
-- +goose StatementEnd
