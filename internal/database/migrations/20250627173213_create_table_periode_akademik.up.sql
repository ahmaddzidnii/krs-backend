CREATE TYPE jenis_semester_enum AS ENUM ('GANJIL', 'GENAP', 'PENDEK');

CREATE TABLE periode_akademik (
    id_periode UUID DEFAULT gen_random_uuid(),
    tahun_akademik VARCHAR(9) NOT NULL UNIQUE,
    jenis_semester jenis_semester_enum NOT NULL,
    tanggal_mulai_krs DATE NOT NULL,
    tanggal_selesai_krs DATE NOT NULL,
    jam_mulai_harian_krs TIME NOT NULL,
    jam_selesai_harian_krs TIME NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_periode)
);


CREATE TRIGGER set_updated_at_periode_akademik
    BEFORE UPDATE ON periode_akademik
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();
