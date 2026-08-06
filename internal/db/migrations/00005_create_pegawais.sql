-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pegawais (
    id                                  BIGSERIAL PRIMARY KEY,
    uuid                                UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id                             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bidang_id                           BIGINT NOT NULL REFERENCES bidangs(id),
    pangkat_golongan_id                 BIGINT REFERENCES pangkat_golongans(id),
    tipe                                VARCHAR(20) NOT NULL DEFAULT 'pegawai'
                                            CONSTRAINT chk_pegawais_tipe CHECK (tipe IN ('pegawai', 'admin')),
    status                              VARCHAR(20) NOT NULL DEFAULT 'aktif'
                                            CONSTRAINT chk_pegawais_status CHECK (status IN ('aktif', 'non aktif')),
    nip                                 VARCHAR(255) NOT NULL UNIQUE,
    nama                                VARCHAR(255) NOT NULL,
    jabatan                             VARCHAR(255),
    kategori_jabatan                    VARCHAR(50)
                                            CONSTRAINT chk_pegawais_kategori_jabatan
                                            CHECK (kategori_jabatan IN ('struktural', 'fungsional auditor', 'fungsional tertentu')),
    peran                               VARCHAR(255),
    kategori_kebutuhan_jam_pelatihan    VARCHAR(20)
                                            CONSTRAINT chk_pegawais_kategori_kebutuhan
                                            CHECK (kategori_kebutuhan_jam_pelatihan IN ('admin', 'pejabat', 'auditor')),
    created_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                          TIMESTAMPTZ,
    created_by                          BIGINT REFERENCES users(id),
    updated_by                          BIGINT REFERENCES users(id),
    deleted_by                          BIGINT REFERENCES users(id)
);

CREATE INDEX idx_pegawais_uuid_hash ON pegawais USING HASH (uuid);
CREATE INDEX idx_pegawais_tipe ON pegawais (tipe);
CREATE INDEX idx_pegawais_status ON pegawais (status);
CREATE INDEX idx_pegawais_bidang_id ON pegawais (bidang_id);
CREATE INDEX idx_pegawais_pangkat_golongan_id ON pegawais (pangkat_golongan_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pegawais;
-- +goose StatementEnd
