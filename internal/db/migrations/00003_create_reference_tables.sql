-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bidangs (
    id         BIGSERIAL PRIMARY KEY,
    uuid       UUID NOT NULL DEFAULT gen_random_uuid(),
    bidang     VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id)
);

CREATE INDEX idx_bidangs_uuid_hash ON bidangs USING HASH (uuid);

CREATE TABLE IF NOT EXISTS pangkat_golongans (
    id         BIGSERIAL PRIMARY KEY,
    uuid       UUID NOT NULL DEFAULT gen_random_uuid(),
    jenjang    VARCHAR(255) NOT NULL,
    pangkat    VARCHAR(255) NOT NULL,
    golongan   VARCHAR(255) NOT NULL,
    ruang      VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    created_by BIGINT REFERENCES users(id),
    updated_by BIGINT REFERENCES users(id),
    deleted_by BIGINT REFERENCES users(id)
);

CREATE INDEX idx_pangkat_golongans_uuid_hash ON pangkat_golongans USING HASH (uuid);

CREATE TABLE IF NOT EXISTS jenis_pelatihans (
    id              BIGSERIAL PRIMARY KEY,
    uuid            UUID NOT NULL DEFAULT gen_random_uuid(),
    jenis_pelatihan VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    created_by      BIGINT REFERENCES users(id),
    updated_by      BIGINT REFERENCES users(id),
    deleted_by      BIGINT REFERENCES users(id)
);

CREATE INDEX idx_jenis_pelatihans_uuid_hash ON jenis_pelatihans USING HASH (uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jenis_pelatihans;
DROP TABLE IF EXISTS pangkat_golongans;
DROP TABLE IF EXISTS bidangs;
-- +goose StatementEnd
