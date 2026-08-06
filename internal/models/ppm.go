package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Ppm struct {
	bun.BaseModel `bun:"table:ppms,alias:ppm"`

	BaseID
	Audit

	PegawaiID          int64     `bun:"pegawai_id,notnull" json:"-"`
	NomorSurat         string    `bun:"nomor_surat,notnull" json:"-"`
	MateriPengembangan string    `bun:"materi_pengembangan,notnull" json:"-"`
	TanggalPelaksanaan time.Time `bun:"tanggal_pelaksanaan,notnull" json:"-"`
	JumlahJamPelatihan int       `bun:"jumlah_jam_pelatihan,notnull" json:"-"`

	Pegawai *Pegawai `bun:"rel:belongs-to,join:pegawai_id=id" json:"-"`
}
