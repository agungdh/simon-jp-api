package models

import "github.com/uptrace/bun"

type JenisPelatihan struct {
	bun.BaseModel `bun:"table:jenis_pelatihans,alias:jp"`

	BaseID
	Audit

	JenisPelatihan string `bun:"jenis_pelatihan,notnull" json:"-"`
}
