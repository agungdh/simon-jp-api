package models

import "github.com/uptrace/bun"

type Pegawai struct {
	bun.BaseModel `bun:"table:pegawais,alias:p"`

	BaseID
	Audit

	UserID              int64   `bun:"user_id,notnull" json:"-"`
	BidangID            int64   `bun:"bidang_id,notnull" json:"-"`
	PangkatGolonganID   *int64  `bun:"pangkat_golongan_id" json:"-"`
	Tipe                string  `bun:"tipe,notnull" json:"-"`
	Status              string  `bun:"status,notnull" json:"-"`
	NIP                 string  `bun:"nip,notnull,unique" json:"-"`
	Nama                string  `bun:"nama,notnull" json:"-"`
	Jabatan             *string `bun:"jabatan" json:"-"`
	KategoriJabatan     *string `bun:"kategori_jabatan" json:"-"`
	Peran               *string `bun:"peran" json:"-"`
	KategoriKebutuhanJP *string `bun:"kategori_kebutuhan_jam_pelatihan" json:"-"`

	User            *User            `bun:"rel:belongs-to,join:user_id=id" json:"-"`
	Bidang          *Bidang          `bun:"rel:belongs-to,join:bidang_id=id" json:"-"`
	PangkatGolongan *PangkatGolongan `bun:"rel:belongs-to,join:pangkat_golongan_id=id" json:"-"`
}
