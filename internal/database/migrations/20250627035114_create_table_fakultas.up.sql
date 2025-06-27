CREATE TABLE fakultas (
    id_fakultas UUID NOT NULL DEFAULT uuid_generate_v4(),
    kode_fakultas VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    singkatan VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_fakultas)
);

CREATE TRIGGER set_updated_at_fakultas
    BEFORE UPDATE ON fakultas
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();