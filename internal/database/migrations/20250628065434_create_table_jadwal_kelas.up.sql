CREATE TABLE jadwal_kelas (
    id_jadwal UUID NOT NULL DEFAULT gen_random_uuid(),
    id_kelas UUID NOT NULL,
    hari VARCHAR(10) NOT NULL,
    waktu_mulai TIME NOT NULL,
    waktu_selesai TIME NOT NULL,
    ruang VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (id_kelas) REFERENCES kelas_ditawarkan(id_kelas) ON DELETE CASCADE,
    PRIMARY KEY (id_jadwal)
);


CREATE TRIGGER set_updated_at_jadwal_kelas
    BEFORE UPDATE ON jadwal_kelas
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();