CREATE TYPE jenjang_enum AS ENUM ('Sarjana (S1)', 'Magister (S2)');

CREATE TABLE program_studi (
    id_prodi UUID NOT NULL DEFAULT uuid_generate_v4(),
    id_fakultas UUID NOT NULL,
    kode_prodi VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    jenjang jenjang_enum NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_prodi)
);

CREATE TRIGGER set_updated_at_program_studi
    BEFORE UPDATE ON program_studi
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();