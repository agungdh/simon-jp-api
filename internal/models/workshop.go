package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Workshop struct {
	bun.BaseModel `bun:"table:workshops,alias:ws"`

	BaseID
	Audit

	PegawaiID          int64     `bun:"pegawai_id,notnull" json:"-"`
	MateriPengembangan string    `bun:"materi_pengembangan,notnull" json:"-"`
	TanggalPelaksanaan time.Time `bun:"tanggal_pelaksanaan,notnull" json:"-"`
	JumlahJam          int       `bun:"jumlah_jam,notnull" json:"-"`
	Filename           *string   `bun:"filename" json:"-"`

	Pegawai *Pegawai `bun:"rel:belongs-to,join:pegawai_id=id" json:"-"`
}
