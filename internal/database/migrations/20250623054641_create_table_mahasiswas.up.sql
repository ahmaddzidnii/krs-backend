CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE status_mahasiswa_enum AS ENUM ('Aktif', 'Cuti', 'Non-Aktif');
CREATE TYPE status_pembayaran_enum AS ENUM ('Lunas', 'Belum Lunas');

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Membuat fungsi yang berisi logika untuk menghitung jatah SKS
-- berdasarkan kolom 'ips_lalu'.
CREATE OR REPLACE FUNCTION hitung_jatah_sks_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.ips_lalu >= 3.00 THEN
        NEW.jatah_sks := 24;
    ELSIF NEW.ips_lalu >= 2.50 THEN
        NEW.jatah_sks := 22;
    ELSIF NEW.ips_lalu >= 2.00 THEN
        NEW.jatah_sks := 20;
    ELSIF NEW.ips_lalu >= 1.50 THEN
        NEW.jatah_sks := 18;
ELSE
        NEW.jatah_sks := 16;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE mahasiswas (
    id_mahasiswa UUID NOT NULL DEFAULT uuid_generate_v4(),
    id_user UUID NOT NULL UNIQUE,
    id_prodi UUID NOT NULL,
    id_dpa UUID NOT NULL,
    id_kurikulum UUID NOT NULL,
    nim VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    ipk NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    ips_lalu NUMERIC(3, 2) NOT NULL DEFAULT 0.00,
    semester_berjalan INT NOT NULL DEFAULT 1,
    sks_kumulatif INT NOT NULL DEFAULT 0,
    jatah_sks INT NOT NULL DEFAULT 0,
    status_mahasiswa status_mahasiswa_enum NOT NULL DEFAULT 'Aktif',
    status_pembayaran status_pembayaran_enum NOT NULL DEFAULT 'Belum Lunas',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_mahasiswa)
);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON mahasiswas
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TRIGGER trigger_hitung_jatah_sks
    BEFORE INSERT OR UPDATE ON mahasiswas
    FOR EACH ROW
    EXECUTE PROCEDURE hitung_jatah_sks_trigger_func();
