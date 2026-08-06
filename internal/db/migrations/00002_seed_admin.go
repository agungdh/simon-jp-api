package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upSeedAdmin, downSeedAdmin)
}

func upSeedAdmin(ctx context.Context, tx *sql.Tx) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING
	`, "admin", string(hash))
	return err
}

func downSeedAdmin(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE username = $1`, "admin")
	return err
}
