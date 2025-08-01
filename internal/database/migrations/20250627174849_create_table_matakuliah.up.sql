CREATE TABLE mata_kuliah (
    id_matkul UUID DEFAULT gen_random_uuid(),
    kode_matkul VARCHAR(256) NOT NULL UNIQUE,
    nama VARCHAR(256) NOT NULL,
    sks INT NOT NULL ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_matkul)
);


CREATE TRIGGER set_updated_at_mata_kuliah
    BEFORE UPDATE ON mata_kuliah
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();
