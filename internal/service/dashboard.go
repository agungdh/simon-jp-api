package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"simon-jp-api/internal/models"
	"simon-jp-api/internal/repository"
)

const dateLayout = "2006-01-02"

type DetailItem struct {
	UUID      string
	Tipe      string
	Tanggal   time.Time
	Materi    string
	JumlahJam int
}

type DetailCapaian struct {
	Diklat  int
	Ppm     int
	Seminar int
	Webinar int
	Lc      int
}

type Rekapitulasi struct {
	Triwulan1 int
	Triwulan2 int
	Triwulan3 int
	Triwulan4 int
}

type DashboardResult struct {
	Pegawai       *models.Pegawai
	Tahun         int
	Kategori      *string
	JumlahMinimal int
	JumlahCapaian int
	DetailCapaian DetailCapaian
	Rekapitulasi  Rekapitulasi
	DetailItems   []DetailItem
}

type DashboardService struct {
	pegawai *repository.PegawaiRepository
	capaian *repository.CapaianRepository
}

func NewDashboardService(pegawai *repository.PegawaiRepository, capaian *repository.CapaianRepository) *DashboardService {
	return &DashboardService{pegawai: pegawai, capaian: capaian}
}

func (s *DashboardService) Get(ctx context.Context, userID int64, tahun int) (*DashboardResult, error) {
	pegawai, err := s.pegawai.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if tahun == 0 {
		tahun = time.Now().Year()
	}

	res := &DashboardResult{
		Pegawai:  pegawai,
		Tahun:    tahun,
		Kategori: pegawai.KategoriKebutuhanJP,
	}
	if pegawai.KategoriKebutuhanJP != nil {
		res.JumlahMinimal = KategoriJumlahMinimal(*pegawai.KategoriKebutuhanJP)
	}

	sums := func(quarter int) int {
		return s.sumAll(ctx, pegawai.ID, tahun, quarter)
	}
	res.DetailCapaian = DetailCapaian{
		Diklat:  s.sumTable(ctx, pegawai.ID, tahun, "diklats", "jumlah_jam_pelatihan", "dari_tanggal_pelaksanaan"),
		Ppm:     s.sumTable(ctx, pegawai.ID, tahun, "ppms", "jumlah_jam_pelatihan", "tanggal_pelaksanaan"),
		Seminar: s.sumTable(ctx, pegawai.ID, tahun, "seminars", "jumlah_jam", "tanggal_pelaksanaan"),
		Webinar: s.sumTable(ctx, pegawai.ID, tahun, "webinars", "jumlah_jam", "tanggal_pelaksanaan"),
		Lc:      s.sumTable(ctx, pegawai.ID, tahun, "lcs", "jumlah_jam", "tanggal_pelaksanaan"),
	}
	res.JumlahCapaian = res.DetailCapaian.Diklat + res.DetailCapaian.Ppm + res.DetailCapaian.Seminar + res.DetailCapaian.Webinar + res.DetailCapaian.Lc
	res.Rekapitulasi = Rekapitulasi{
		Triwulan1: sums(1),
		Triwulan2: sums(2),
		Triwulan3: sums(3),
		Triwulan4: sums(4),
	}

	items, err := s.detailItems(ctx, pegawai.ID, tahun)
	if err != nil {
		return nil, err
	}
	res.DetailItems = items

	return res, nil
}

func (s *DashboardService) sumTable(ctx context.Context, pegawaiID int64, tahun int, table, jamColumn, dateColumn string) int {
	sum, err := s.capaian.SumByYear(ctx, pegawaiID, tahun, table, jamColumn, dateColumn)
	if err != nil {
		return 0
	}
	return sum
}

func (s *DashboardService) sumAll(ctx context.Context, pegawaiID int64, tahun, quarter int) int {
	sum := 0
	sum += s.sumQuarter(ctx, pegawaiID, tahun, quarter, "diklats", "jumlah_jam_pelatihan", "dari_tanggal_pelaksanaan")
	sum += s.sumQuarter(ctx, pegawaiID, tahun, quarter, "ppms", "jumlah_jam_pelatihan", "tanggal_pelaksanaan")
	sum += s.sumQuarter(ctx, pegawaiID, tahun, quarter, "seminars", "jumlah_jam", "tanggal_pelaksanaan")
	sum += s.sumQuarter(ctx, pegawaiID, tahun, quarter, "webinars", "jumlah_jam", "tanggal_pelaksanaan")
	sum += s.sumQuarter(ctx, pegawaiID, tahun, quarter, "lcs", "jumlah_jam", "tanggal_pelaksanaan")
	return sum
}

func (s *DashboardService) sumQuarter(ctx context.Context, pegawaiID int64, tahun, quarter int, table, jamColumn, dateColumn string) int {
	sum, err := s.capaian.SumByQuarter(ctx, pegawaiID, tahun, quarter, table, jamColumn, dateColumn)
	if err != nil {
		return 0
	}
	return sum
}

func (s *DashboardService) detailItems(ctx context.Context, pegawaiID int64, tahun int) ([]DetailItem, error) {
	items := make([]DetailItem, 0)

	diklats, err := s.capaian.ListDiklatsByYear(ctx, pegawaiID, tahun)
	if err != nil {
		return nil, err
	}
	for _, d := range diklats {
		items = append(items, DetailItem{
			UUID:      d.UUID,
			Tipe:      "Diklat",
			Tanggal:   d.DariTanggalPelaksanaan,
			Materi:    d.MateriPengembangan,
			JumlahJam: d.JumlahJamPelatihan,
		})
	}

	ppms, err := s.capaian.ListPpmsByYear(ctx, pegawaiID, tahun)
	if err != nil {
		return nil, err
	}
	for _, p := range ppms {
		items = append(items, DetailItem{
			UUID:      p.UUID,
			Tipe:      "PPM",
			Tanggal:   p.TanggalPelaksanaan,
			Materi:    p.MateriPengembangan,
			JumlahJam: p.JumlahJamPelatihan,
		})
	}

	seminars, err := s.capaian.ListSeminarsByYear(ctx, pegawaiID, tahun)
	if err != nil {
		return nil, err
	}
	for _, v := range seminars {
		items = append(items, DetailItem{
			UUID:      v.UUID,
			Tipe:      "Seminar",
			Tanggal:   v.TanggalPelaksanaan,
			Materi:    v.MateriPengembangan,
			JumlahJam: v.JumlahJam,
		})
	}

	webinars, err := s.capaian.ListWebinarsByYear(ctx, pegawaiID, tahun)
	if err != nil {
		return nil, err
	}
	for _, v := range webinars {
		items = append(items, DetailItem{
			UUID:      v.UUID,
			Tipe:      "Webinar",
			Tanggal:   v.TanggalPelaksanaan,
			Materi:    v.MateriPengembangan,
			JumlahJam: v.JumlahJam,
		})
	}

	lcs, err := s.capaian.ListLcsByYear(ctx, pegawaiID, tahun)
	if err != nil {
		return nil, err
	}
	for _, v := range lcs {
		items = append(items, DetailItem{
			UUID:      v.UUID,
			Tipe:      "LC",
			Tanggal:   v.TanggalPelaksanaan,
			Materi:    v.MateriPengembangan,
			JumlahJam: v.JumlahJam,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Tanggal.Before(items[j].Tanggal)
	})

	return items, nil
}
