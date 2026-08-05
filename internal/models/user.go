package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int64     `bun:"id,pk,autoincrement" json:"id"`
	Username     string    `bun:"username,notnull,unique" json:"username"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	CreatedAt    time.Time `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,notnull" json:"updated_at"`
}
