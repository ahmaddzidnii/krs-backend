CREATE TYPE jenis_mata_kuliah AS ENUM (
    'Wajib',
    'Pilihan'
);

CREATE TABLE detail_kurikulum (
    id_detail_kurikulum UUID NOT NULL DEFAULT gen_random_uuid(),
    id_kurikulum UUID NOT NULL,
    id_matkul UUID NOT NULL,
    jenis_matkul jenis_mata_kuliah NOT NULL,
    semester_paket INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (id_kurikulum) REFERENCES kurikulum(id_kurikulum) ON DELETE CASCADE,
    FOREIGN KEY (id_matkul) REFERENCES mata_kuliah(id_matkul) ON DELETE CASCADE,

    PRIMARY KEY (id_matkul)
);

CREATE TRIGGER set_updated_at_detail_kurikulum
    BEFORE UPDATE ON detail_kurikulum
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();