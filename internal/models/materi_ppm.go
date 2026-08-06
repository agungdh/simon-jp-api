package models

import (
	"time"

	"github.com/uptrace/bun"
)

type MateriPpm struct {
	bun.BaseModel `bun:"table:materi_ppms,alias:mp"`

	BaseID
	Audit

	NomorSurat         *string   `bun:"nomor_surat" json:"-"`
	NamaPemateri       *string   `bun:"nama_pemateri" json:"-"`
	MateriPengembangan string    `bun:"materi_pengembangan,notnull" json:"-"`
	TanggalPelaksanaan time.Time `bun:"tanggal_pelaksanaan,notnull" json:"-"`
	LinkMateri         *string   `bun:"link_materi" json:"-"`
	LinkDokumentasi    *string   `bun:"link_dokumentasi" json:"-"`
}
