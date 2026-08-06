package httpapi

import (
	"time"

	"simon-jp-api/internal/models"
)

const dateLayout = "2006-01-02"

type userResponse struct {
	UUID     string       `json:"uuid"`
	Username string       `json:"username"`
	Name     *string      `json:"name"`
	Role     string       `json:"role"`
	Profile  *pegawaiShow `json:"profile"`
}

func newUserResponse(u *models.User, p *models.Pegawai) userResponse {
	res := userResponse{
		UUID:     u.UUID,
		Username: u.Username,
		Role:     u.Role,
	}
	if p != nil {
		name := p.Nama
		res.Name = &name
		profile := newPegawaiShow(p)
		res.Profile = &profile
	}
	return res
}

type pegawaiBrief struct {
	UUID string `json:"uuid"`
	NIP  string `json:"nip"`
	Nama string `json:"nama"`
}

func newPegawaiBrief(p *models.Pegawai) pegawaiBrief {
	return pegawaiBrief{UUID: p.UUID, NIP: p.NIP, Nama: p.Nama}
}

type bidangDTO struct {
	UUID   string `json:"uuid"`
	Bidang string `json:"bidang"`
}

func newBidangDTO(b *models.Bidang) *bidangDTO {
	if b == nil {
		return nil
	}
	return &bidangDTO{UUID: b.UUID, Bidang: b.Bidang}
}

type pangkatDTO struct {
	UUID     string `json:"uuid"`
	Jenjang  string `json:"jenjang"`
	Pangkat  string `json:"pangkat"`
	Golongan string `json:"golongan"`
	Ruang    string `json:"ruang"`
}

func newPangkatDTO(p *models.PangkatGolongan) *pangkatDTO {
	if p == nil {
		return nil
	}
	return &pangkatDTO{UUID: p.UUID, Jenjang: p.Jenjang, Pangkat: p.Pangkat, Golongan: p.Golongan, Ruang: p.Ruang}
}

type jenisPelatihanDTO struct {
	UUID           string `json:"uuid"`
	JenisPelatihan string `json:"jenis_pelatihan"`
}

func newJenisPelatihanDTO(j *models.JenisPelatihan) *jenisPelatihanDTO {
	if j == nil {
		return nil
	}
	return &jenisPelatihanDTO{UUID: j.UUID, JenisPelatihan: j.JenisPelatihan}
}

type pegawaiShow struct {
	UUID                          string      `json:"uuid"`
	NIP                           string      `json:"nip"`
	Nama                          string      `json:"nama"`
	Jabatan                       *string     `json:"jabatan"`
	Peran                         *string     `json:"peran"`
	Status                        string      `json:"status"`
	Tipe                          string      `json:"tipe"`
	KategoriJabatan               *string     `json:"kategori_jabatan"`
	KategoriKebutuhanJamPelatihan *string     `json:"kategori_kebutuhan_jam_pelatihan"`
	Bidang                        *bidangDTO  `json:"bidang"`
	PangkatGolongan               *pangkatDTO `json:"pangkat_golongan"`
	CreatedAt                     time.Time   `json:"created_at"`
	UpdatedAt                     time.Time   `json:"updated_at"`
}

func newPegawaiShow(p *models.Pegawai) pegawaiShow {
	return pegawaiShow{
		UUID:                          p.UUID,
		NIP:                           p.NIP,
		Nama:                          p.Nama,
		Jabatan:                       p.Jabatan,
		Peran:                         p.Peran,
		Status:                        p.Status,
		Tipe:                          p.Tipe,
		KategoriJabatan:               p.KategoriJabatan,
		KategoriKebutuhanJamPelatihan: p.KategoriKebutuhanJP,
		Bidang:                        newBidangDTO(p.Bidang),
		PangkatGolongan:               newPangkatDTO(p.PangkatGolongan),
		CreatedAt:                     p.CreatedAt,
		UpdatedAt:                     p.UpdatedAt,
	}
}

type pegawaiListItem struct {
	UUID            string  `json:"uuid"`
	NIP             string  `json:"nip"`
	Nama            string  `json:"nama"`
	Jabatan         *string `json:"jabatan"`
	Status          string  `json:"status"`
	Bidang          *string `json:"bidang"`
	PangkatGolongan *string `json:"pangkat_golongan"`
}

func newPegawaiListItem(p *models.Pegawai) pegawaiListItem {
	item := pegawaiListItem{
		UUID:    p.UUID,
		NIP:     p.NIP,
		Nama:    p.Nama,
		Jabatan: p.Jabatan,
		Status:  p.Status,
	}
	if p.Bidang != nil {
		item.Bidang = &p.Bidang.Bidang
	}
	if p.PangkatGolongan != nil {
		pg := p.PangkatGolongan.Pangkat + "/" + p.PangkatGolongan.Golongan
		item.PangkatGolongan = &pg
	}
	return item
}

type diklatListItem struct {
	UUID                     string             `json:"uuid"`
	Pegawai                  pegawaiBrief       `json:"pegawai"`
	JenisPelatihan           *jenisPelatihanDTO `json:"jenis_pelatihan"`
	NomorSurat               string             `json:"nomor_surat"`
	MateriPengembangan       string             `json:"materi_pengembangan"`
	DariTanggalPelaksanaan   string             `json:"dari_tanggal_pelaksanaan"`
	SampaiTanggalPelaksanaan string             `json:"sampai_tanggal_pelaksanaan"`
	JumlahJamPelatihan       int                `json:"jumlah_jam_pelatihan"`
}

func newDiklatListItem(d *models.Diklat) diklatListItem {
	return diklatListItem{
		UUID:                     d.UUID,
		Pegawai:                  newPegawaiBrief(d.Pegawai),
		JenisPelatihan:           newJenisPelatihanDTO(d.JenisPelatihan),
		NomorSurat:               d.NomorSurat,
		MateriPengembangan:       d.MateriPengembangan,
		DariTanggalPelaksanaan:   d.DariTanggalPelaksanaan.Format(dateLayout),
		SampaiTanggalPelaksanaan: d.SampaiTanggalPelaksanaan.Format(dateLayout),
		JumlahJamPelatihan:       d.JumlahJamPelatihan,
	}
}

type diklatShow struct {
	diklatListItem
	Filename  *string   `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newDiklatShow(d *models.Diklat) diklatShow {
	return diklatShow{
		diklatListItem: newDiklatListItem(d),
		Filename:       d.Filename,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

type ppmListItem struct {
	UUID               string       `json:"uuid"`
	Pegawai            pegawaiBrief `json:"pegawai"`
	NomorSurat         string       `json:"nomor_surat"`
	MateriPengembangan string       `json:"materi_pengembangan"`
	TanggalPelaksanaan string       `json:"tanggal_pelaksanaan"`
	JumlahJamPelatihan int          `json:"jumlah_jam_pelatihan"`
}

func newPpmListItem(p *models.Ppm) ppmListItem {
	return ppmListItem{
		UUID:               p.UUID,
		Pegawai:            newPegawaiBrief(p.Pegawai),
		NomorSurat:         p.NomorSurat,
		MateriPengembangan: p.MateriPengembangan,
		TanggalPelaksanaan: p.TanggalPelaksanaan.Format(dateLayout),
		JumlahJamPelatihan: p.JumlahJamPelatihan,
	}
}

type ppmShow struct {
	ppmListItem
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newPpmShow(p *models.Ppm) ppmShow {
	return ppmShow{
		ppmListItem: newPpmListItem(p),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type activityListItem struct {
	UUID               string       `json:"uuid"`
	Pegawai            pegawaiBrief `json:"pegawai"`
	MateriPengembangan string       `json:"materi_pengembangan"`
	TanggalPelaksanaan string       `json:"tanggal_pelaksanaan"`
	JumlahJam          int          `json:"jumlah_jam"`
}

type activityShow struct {
	activityListItem
	Filename  *string   `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type materiPpmListItem struct {
	UUID               string  `json:"uuid"`
	NomorSurat         *string `json:"nomor_surat"`
	NamaPemateri       *string `json:"nama_pemateri"`
	MateriPengembangan string  `json:"materi_pengembangan"`
	TanggalPelaksanaan string  `json:"tanggal_pelaksanaan"`
	LinkMateri         *string `json:"link_materi"`
	LinkDokumentasi    *string `json:"link_dokumentasi"`
}

func newMateriPpmListItem(m *models.MateriPpm) materiPpmListItem {
	return materiPpmListItem{
		UUID:               m.UUID,
		NomorSurat:         m.NomorSurat,
		NamaPemateri:       m.NamaPemateri,
		MateriPengembangan: m.MateriPengembangan,
		TanggalPelaksanaan: m.TanggalPelaksanaan.Format(dateLayout),
		LinkMateri:         m.LinkMateri,
		LinkDokumentasi:    m.LinkDokumentasi,
	}
}

type materiPpmShow struct {
	materiPpmListItem
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newMateriPpmShow(m *models.MateriPpm) materiPpmShow {
	return materiPpmShow{
		materiPpmListItem: newMateriPpmListItem(m),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}
