package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upSeedDemo, downSeedDemo)
}

func queryID(ctx context.Context, tx *sql.Tx, query string, args ...any) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, query, args...).Scan(&id)
	return id, err
}

func upSeedDemo(ctx context.Context, tx *sql.Tx) error {
	year := time.Now().Year()

	demoHash, err := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	adminUserID, err := queryID(ctx, tx, `SELECT id FROM users WHERE username = 'admin' AND deleted_at IS NULL`)
	if err != nil {
		return err
	}

	demoUserID, err := queryID(ctx, tx, `SELECT id FROM users WHERE username = 'demo' AND deleted_at IS NULL`)
	if err != nil {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO users (username, password, role)
			VALUES ('demo', $1, 'pegawai')
		`, string(demoHash)); err != nil {
			return err
		}
		demoUserID, err = queryID(ctx, tx, `SELECT id FROM users WHERE username = 'demo'`)
		if err != nil {
			return err
		}
	}

	adminBidangID, err := queryID(ctx, tx, `SELECT id FROM bidangs WHERE bidang = 'SUBBAGIAN KEPEGAWAIAN'`)
	if err != nil {
		return err
	}

	demoBidangID, err := queryID(ctx, tx, `SELECT id FROM bidangs WHERE bidang = 'BIDANG PENGAWASAN INSTANSI PEMERINTAH PUSAT'`)
	if err != nil {
		return err
	}

	pangkatID, err := queryID(ctx, tx, `SELECT id FROM pangkat_golongans WHERE pangkat = 'Pembina' AND golongan = 'IV' AND ruang = 'a'`)
	if err != nil {
		return err
	}

	jenisPenjenjanganID, err := queryID(ctx, tx, `SELECT id FROM jenis_pelatihans WHERE jenis_pelatihan = 'Penjenjangan'`)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO pegawais (user_id, bidang_id, tipe, status, nip, nama, jabatan, kategori_jabatan, peran, kategori_kebutuhan_jam_pelatihan)
		VALUES ($1, $2, 'admin', 'aktif', '197001012010011001', 'ADMIN SISTEM', 'Kepala Subbagian Kepegawaian', 'struktural', 'Admin', 'admin')
		ON CONFLICT (nip) DO NOTHING
	`, adminUserID, adminBidangID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO pegawais (user_id, bidang_id, pangkat_golongan_id, tipe, status, nip, nama, jabatan, kategori_jabatan, peran, kategori_kebutuhan_jam_pelatihan)
		VALUES ($1, $2, $3, 'pegawai', 'aktif', '198801012019031001', 'DEMO PEGAWAI', 'Auditor Madya', 'fungsional auditor', 'Auditor', 'auditor')
		ON CONFLICT (nip) DO NOTHING
	`, demoUserID, demoBidangID, pangkatID); err != nil {
		return err
	}

	demoPegawaiID, err := queryID(ctx, tx, `SELECT id FROM pegawais WHERE user_id = $1`, demoUserID)
	if err != nil {
		return err
	}

	diklat := []struct {
		month  time.Month
		nomor  string
		materi string
		jam    int
		days   int
	}{
		{time.January, "Diklat-20260001", "Audit Kinerja Berbasis Risiko", 20, 3},
		{time.July, "Diklat-20260002", "Teknik Pemeriksaan Investigatif", 16, 2},
	}
	for _, d := range diklat {
		start := time.Date(year, d.month, 10, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 0, d.days)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO diklats (pegawai_id, jenis_pelatihan_id, nomor_surat, materi_pengembangan, dari_tanggal_pelaksanaan, sampai_tanggal_pelaksanaan, jumlah_jam_pelatihan)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, demoPegawaiID, jenisPenjenjanganID, d.nomor, d.materi, start, end, d.jam); err != nil {
			return err
		}
	}

	ppm := []struct {
		month  time.Month
		nomor  string
		materi string
		jam    int
	}{
		{time.February, "PPM-20260001", "Reviu Laporan Keuangan", 8},
		{time.May, "PPM-20260002", "Audit Kinerja Program Prioritas", 8},
		{time.October, "PPM-20260003", "Evaluasi Capaian Output", 6},
	}
	for i, p := range ppm {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO ppms (pegawai_id, nomor_surat, materi_pengembangan, tanggal_pelaksanaan, jumlah_jam_pelatihan)
			VALUES ($1, $2, $3, $4, $5)
		`, demoPegawaiID, p.nomor, p.materi, time.Date(year, p.month, 15+i, 0, 0, 0, 0, time.UTC), p.jam); err != nil {
			return err
		}
	}

	insertActivity := func(table string, month time.Month, materi string, jam int) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO `+table+` (pegawai_id, materi_pengembangan, tanggal_pelaksanaan, jumlah_jam)
			VALUES ($1, $2, $3, $4)
		`, demoPegawaiID, materi, time.Date(year, month, 20, 0, 0, 0, 0, time.UTC), jam)
		return err
	}

	if err := insertActivity("seminars", time.March, "Sosialisasi Standar Audit", 4); err != nil {
		return err
	}
	if err := insertActivity("seminars", time.September, "Seminar Nasional Pengawasan", 4); err != nil {
		return err
	}
	if err := insertActivity("webinars", time.April, "Webinar Fraud Risk Assessment", 2); err != nil {
		return err
	}
	if err := insertActivity("webinars", time.August, "Webinar Audit Berbasis Data", 2); err != nil {
		return err
	}
	if err := insertActivity("lcs", time.June, "Learning Center JFA", 2); err != nil {
		return err
	}

	return nil
}

func downSeedDemo(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM pegawais WHERE nip IN ('197001012010011001', '198801012019031001')
	`); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE username = 'demo'`)
	return err
}
