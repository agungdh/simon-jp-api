-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ppms (
    id                             BIGSERIAL PRIMARY KEY,
    uuid                           UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id                     BIGINT NOT NULL REFERENCES pegawais(id),
    nomor_surat                    VARCHAR(255) NOT NULL,
    materi_pengembangan            VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan            DATE NOT NULL,
    jumlah_jam_pelatihan           INTEGER NOT NULL,
    created_at                     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                     TIMESTAMPTZ,
    created_by                     BIGINT REFERENCES users(id),
    updated_by                     BIGINT REFERENCES users(id),
    deleted_by                     BIGINT REFERENCES users(id)
);

CREATE INDEX idx_ppms_uuid_hash ON ppms USING HASH (uuid);
CREATE INDEX idx_ppms_nomor_surat ON ppms (nomor_surat);
CREATE INDEX idx_ppms_tanggal ON ppms (tanggal_pelaksanaan);
CREATE INDEX idx_ppms_pegawai_id ON ppms (pegawai_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ppms;
-- +goose StatementEnd
