CREATE OR REPLACE FUNCTION update_krs_summary()
RETURNS TRIGGER AS $$
DECLARE
    sks_to_change INT;
    target_krs_id UUID;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        target_krs_id := NEW.id_krs;
    ELSIF (TG_OP = 'DELETE') THEN
        target_krs_id := OLD.id_krs;
    END IF;

    SELECT mk.sks INTO sks_to_change
    FROM kelas_ditawarkan kd
             JOIN mata_kuliah mk ON kd.id_matkul = mk.id_matkul
    WHERE kd.id_kelas = CASE
                            WHEN TG_OP = 'INSERT' THEN NEW.id_kelas
                            WHEN TG_OP = 'DELETE' THEN OLD.id_kelas
        END;

    IF (TG_OP = 'INSERT') THEN
        UPDATE krs
        SET
            total_sks_diambil = total_sks_diambil + sks_to_change,
            updated_at = NOW()
        WHERE id_krs = target_krs_id;
        RETURN NEW;

    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE krs
        SET
            total_sks_diambil = total_sks_diambil - sks_to_change,
            updated_at = NOW()
        WHERE id_krs = target_krs_id;
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE krs (
    id_krs UUID NOT NULL DEFAULT gen_random_uuid(),
    id_mahasiswa UUID NOT NULL,
    id_periode UUID NOT NULL,
    total_sks_diambil INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id_krs)
);

CREATE TABLE detail_krs (
    id_krs UUID NOT NULL,
    id_kelas UUID NOT NULL,

    FOREIGN KEY  (id_krs) REFERENCES krs(id_krs) ON DELETE CASCADE,
    FOREIGN KEY  (id_kelas) REFERENCES kelas_ditawarkan(id_kelas) ON DELETE CASCADE,

    PRIMARY KEY (id_krs, id_kelas)
);

CREATE TRIGGER set_updated_at_krs
    BEFORE UPDATE ON krs
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated_at();

CREATE TRIGGER on_detail_krs_insert
    AFTER INSERT ON detail_krs
    FOR EACH ROW
    EXECUTE FUNCTION update_krs_summary();

CREATE TRIGGER on_detail_krs_delete
    AFTER DELETE ON detail_krs
    FOR EACH ROW
    EXECUTE FUNCTION update_krs_summary();