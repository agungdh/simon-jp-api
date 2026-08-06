package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Diklat struct {
	bun.BaseModel `bun:"table:diklats,alias:d"`

	BaseID
	Audit

	PegawaiID                int64     `bun:"pegawai_id,notnull" json:"-"`
	JenisPelatihanID         int64     `bun:"jenis_pelatihan_id,notnull" json:"-"`
	NomorSurat               string    `bun:"nomor_surat,notnull" json:"-"`
	MateriPengembangan       string    `bun:"materi_pengembangan,notnull" json:"-"`
	DariTanggalPelaksanaan   time.Time `bun:"dari_tanggal_pelaksanaan,notnull" json:"-"`
	SampaiTanggalPelaksanaan time.Time `bun:"sampai_tanggal_pelaksanaan,notnull" json:"-"`
	JumlahJamPelatihan       int       `bun:"jumlah_jam_pelatihan,notnull" json:"-"`
	Filename                 *string   `bun:"filename" json:"-"`

	Pegawai        *Pegawai        `bun:"rel:belongs-to,join:pegawai_id=id" json:"-"`
	JenisPelatihan *JenisPelatihan `bun:"rel:belongs-to,join:jenis_pelatihan_id=id" json:"-"`
}
