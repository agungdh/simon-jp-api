package models

import "time"

type BaseID struct {
	ID   int64  `bun:"id,pk,autoincrement" json:"-"`
	UUID string `bun:"uuid,notnull" json:"uuid"`
}

type Audit struct {
	CreatedAt time.Time  `bun:"created_at,notnull" json:"-"`
	UpdatedAt time.Time  `bun:"updated_at,notnull" json:"-"`
	DeletedAt *time.Time `bun:"deleted_at" json:"-"`
	CreatedBy *int64     `bun:"created_by" json:"-"`
	UpdatedBy *int64     `bun:"updated_by" json:"-"`
	DeletedBy *int64     `bun:"deleted_by" json:"-"`
}
