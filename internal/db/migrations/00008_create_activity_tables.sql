-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS seminars (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_seminars_uuid_hash ON seminars USING HASH (uuid);
CREATE INDEX idx_seminars_tanggal ON seminars (tanggal_pelaksanaan);
CREATE INDEX idx_seminars_pegawai_id ON seminars (pegawai_id);

CREATE TABLE IF NOT EXISTS webinars (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_webinars_uuid_hash ON webinars USING HASH (uuid);
CREATE INDEX idx_webinars_tanggal ON webinars (tanggal_pelaksanaan);
CREATE INDEX idx_webinars_pegawai_id ON webinars (pegawai_id);

CREATE TABLE IF NOT EXISTS lcs (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_lcs_uuid_hash ON lcs USING HASH (uuid);
CREATE INDEX idx_lcs_tanggal ON lcs (tanggal_pelaksanaan);
CREATE INDEX idx_lcs_pegawai_id ON lcs (pegawai_id);

CREATE TABLE IF NOT EXISTS belajar_mandiris (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_belajar_mandiris_uuid_hash ON belajar_mandiris USING HASH (uuid);
CREATE INDEX idx_belajar_mandiris_tanggal ON belajar_mandiris (tanggal_pelaksanaan);
CREATE INDEX idx_belajar_mandiris_pegawai_id ON belajar_mandiris (pegawai_id);

CREATE TABLE IF NOT EXISTS mentorings (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_mentorings_uuid_hash ON mentorings USING HASH (uuid);
CREATE INDEX idx_mentorings_tanggal ON mentorings (tanggal_pelaksanaan);
CREATE INDEX idx_mentorings_pegawai_id ON mentorings (pegawai_id);

CREATE TABLE IF NOT EXISTS coachings (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_coachings_uuid_hash ON coachings USING HASH (uuid);
CREATE INDEX idx_coachings_tanggal ON coachings (tanggal_pelaksanaan);
CREATE INDEX idx_coachings_pegawai_id ON coachings (pegawai_id);

CREATE TABLE IF NOT EXISTS workshops (
    id                    BIGSERIAL PRIMARY KEY,
    uuid                  UUID NOT NULL DEFAULT gen_random_uuid(),
    pegawai_id            BIGINT NOT NULL REFERENCES pegawais(id),
    materi_pengembangan   VARCHAR(255) NOT NULL,
    tanggal_pelaksanaan   DATE NOT NULL,
    jumlah_jam            INTEGER NOT NULL,
    filename              VARCHAR(255),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at            TIMESTAMPTZ,
    created_by            BIGINT REFERENCES users(id),
    updated_by            BIGINT REFERENCES users(id),
    deleted_by            BIGINT REFERENCES users(id)
);

CREATE INDEX idx_workshops_uuid_hash ON workshops USING HASH (uuid);
CREATE INDEX idx_workshops_tanggal ON workshops (tanggal_pelaksanaan);
CREATE INDEX idx_workshops_pegawai_id ON workshops (pegawai_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workshops;
DROP TABLE IF EXISTS coachings;
DROP TABLE IF EXISTS mentorings;
DROP TABLE IF EXISTS belajar_mandiris;
DROP TABLE IF EXISTS lcs;
DROP TABLE IF EXISTS webinars;
DROP TABLE IF EXISTS seminars;
-- +goose StatementEnd
