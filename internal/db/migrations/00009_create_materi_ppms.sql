-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS materi_ppms (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    nomor_surat           VARCHAR(255),
    nama_pemateri         VARCHAR(255),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    link_materi           VARCHAR(255),
    link_dokumentasi      VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_materi_ppms_uuid_hash ON materi_ppms USING HASH (uuid);
CREATE INDEX idx_materi_ppms_nomor_surat ON materi_ppms (nomor_surat);
CREATE INDEX idx_materi_ppms_tanggal ON materi_ppms (tanggal_pelaksanaan);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS materi_ppms;
-- +goose StatementEnd
