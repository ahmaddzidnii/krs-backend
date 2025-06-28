CREATE TABLE kurikulum (
    id_kurikulum UUID NOT NULL DEFAULT uuid_generate_v4(),
    id_prodi UUID NOT NULL,
    kode_kurikulum VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    isActive BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_kurikulum)
);

CREATE TRIGGER set_updated_at_kurikulum
    BEFORE UPDATE ON kurikulum
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();