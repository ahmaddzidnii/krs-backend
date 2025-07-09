ALTER TABLE dosen_pengajar_kelas
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE dosen_pengajar_kelas
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
