package models

import "github.com/uptrace/bun"

type Bidang struct {
	bun.BaseModel `bun:"table:bidangs,alias:b"`

	BaseID
	Audit

	Bidang string `bun:"bidang,notnull" json:"-"`
}
