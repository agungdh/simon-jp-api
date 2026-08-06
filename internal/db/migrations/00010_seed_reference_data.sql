-- +goose Up
-- +goose StatementBegin
INSERT INTO bidangs (bidang) VALUES
    ('KEPALA PERWAKILAN'),
    ('KEPALA BAGIAN UMUM'),
    ('SUBBAGIAN KEPEGAWAIAN'),
    ('SUBBAGIAN KEUANGAN'),
    ('SUBBAGIAN UMUM'),
    ('BIDANG PENGAWASAN INSTANSI PEMERINTAH PUSAT'),
    ('BIDANG AKUNTABILITAS PEMERINTAH DAERAH'),
    ('BIDANG AKUNTAN NEGARA'),
    ('BIDANG INVESTIGASI'),
    ('BIDANG PROGRAM DAN PELAPORAN SERTA PEMBINAAN JFA')
ON CONFLICT DO NOTHING;

INSERT INTO jenis_pelatihans (jenis_pelatihan) VALUES
    ('Penjenjangan'),
    ('MOOC'),
    ('Sertifikasi'),
    ('Teknis Substansif'),
    ('Pelatihan Kepemimpinan')
ON CONFLICT DO NOTHING;

INSERT INTO pangkat_golongans (jenjang, pangkat, golongan, ruang) VALUES
    ('Juru', 'Juru Muda', 'I', 'a'),
    ('Juru', 'Juru Muda Tingkat I', 'I', 'b'),
    ('Juru', 'Juru', 'I', 'c'),
    ('Juru', 'Juru Tingkat I', 'I', 'd'),
    ('Pengatur', 'Pengatur Muda', 'II', 'a'),
    ('Pengatur', 'Pengatur Muda Tingkat I', 'II', 'b'),
    ('Pengatur', 'Pengatur', 'II', 'c'),
    ('Pengatur', 'Pengatur Tingkat I', 'II', 'd'),
    ('Penata', 'Penata Muda', 'III', 'a'),
    ('Penata', 'Penata Muda Tingkat I', 'III', 'b'),
    ('Penata', 'Penata', 'III', 'c'),
    ('Penata', 'Penata Tingkat I', 'III', 'd'),
    ('Pembina', 'Pembina', 'IV', 'a'),
    ('Pembina', 'Pembina Tingkat I', 'IV', 'b'),
    ('Pembina Utama', 'Pembina Utama Muda', 'IV', 'c'),
    ('Pembina Utama', 'Pembina Utama Madya', 'IV', 'd'),
    ('Pembina Utama', 'Pembina Utama', 'IV', 'e')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE pangkat_golongans, jenis_pelatihans, bidangs RESTART IDENTITY CASCADE;
-- +goose StatementEnd
