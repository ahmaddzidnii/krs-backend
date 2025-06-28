CREATE TABLE  dosen (
    id_dosen UUID DEFAULT gen_random_uuid(),
    nip VARCHAR(18) NOT NULL UNIQUE,
    id_user UUID NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (id_user) REFERENCES users(id_user) ON DELETE CASCADE,

    PRIMARY KEY (id_dosen)
);

CREATE TABLE  pegawai (
    id_pegawai UUID DEFAULT gen_random_uuid(),
    nip VARCHAR(18) NOT NULL UNIQUE,
    id_user UUID NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (id_user) REFERENCES users(id_user) ON DELETE CASCADE,

    PRIMARY KEY (id_pegawai)
);

CREATE TRIGGER set_updated_at_dosen
    BEFORE UPDATE ON dosen
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TRIGGER set_updated_at_pegawai
    BEFORE UPDATE ON pegawai
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();