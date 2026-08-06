package models

import "github.com/uptrace/bun"

type PangkatGolongan struct {
	bun.BaseModel `bun:"table:pangkat_golongans,alias:pg"`

	BaseID
	Audit

	Jenjang  string `bun:"jenjang,notnull" json:"-"`
	Pangkat  string `bun:"pangkat,notnull" json:"-"`
	Golongan string `bun:"golongan,notnull" json:"-"`
	Ruang    string `bun:"ruang,notnull" json:"-"`
}
