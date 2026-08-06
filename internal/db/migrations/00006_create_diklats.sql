-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS diklats (
    id                             BIGSERIAL PRIMARY KEY,
    uuid                           UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id                     BIGINT NOT NULL REFERENCES pegawais(id),
    jenis_pelatihan_id             BIGINT NOT NULL REFERENCES jenis_pelatihans(id),
    nomor_surat                    VARCHAR(255) NOT NULL,
    materi_pengembangan            VARCHAR(255) NOT NULL,
    dari_tanggal_pelaksanaan       DATE NOT NULL,
    sampai_tanggal_pelaksanaan     DATE NOT NULL,
    jumlah_jam_pelatihan           INTEGER NOT NULL,
    filename                       VARCHAR(255),
    created_at                     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                     TIMESTAMPTZ,
    created_by                     BIGINT REFERENCES users(id),
    updated_by                     BIGINT REFERENCES users(id),
    deleted_by                     BIGINT REFERENCES users(id)
);

CREATE INDEX idx_diklats_uuid_hash ON diklats USING HASH (uuid);
CREATE INDEX idx_diklats_nomor_surat ON diklats (nomor_surat);
CREATE INDEX idx_diklats_dari_tanggal ON diklats (dari_tanggal_pelaksanaan);
CREATE INDEX idx_diklats_sampai_tanggal ON diklats (sampai_tanggal_pelaksanaan);
CREATE INDEX idx_diklats_pegawai_id ON diklats (pegawai_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS diklats;
-- +goose StatementEnd
