CREATE TABLE kelas_ditawarkan (
    id_kelas UUID NOT NULL DEFAULT uuid_generate_v4(),
    id_periode UUID NOT NULL,
    id_matkul UUID NOT NULL,
    kouta INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_kelas)
);

CREATE TABLE dosen_pengajar_kelas (
    id_dosen UUID NOT NULL,
    id_kelas UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_dosen, id_kelas),
    FOREIGN KEY (id_dosen) REFERENCES dosen(id_dosen) ON DELETE CASCADE,
    FOREIGN KEY (id_kelas) REFERENCES kelas_ditawarkan(id_kelas) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at_kelas_ditawarkan
    BEFORE UPDATE ON kelas_ditawarkan
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();