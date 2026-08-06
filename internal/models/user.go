package models

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	BaseID
	Audit

	Username string `bun:"username,notnull,unique" json:"username"`
	Password string `bun:"password,notnull" json:"-"`
	Role     string `bun:"role,notnull" json:"-"`
}
