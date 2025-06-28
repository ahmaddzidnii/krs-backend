ALTER TABLE program_studi
ADD CONSTRAINT fk_fakultas_program_studi_id_fakultas
    FOREIGN KEY (id_fakultas)
    REFERENCES fakultas(id_fakultas)
    ON DELETE CASCADE;

ALTER TABLE mahasiswa
ADD CONSTRAINT fk_mahasiswa_program_studi_id_prodi
    FOREIGN KEY (id_prodi)
    REFERENCES program_studi(id_prodi)
    ON DELETE CASCADE;

ALTER TABLE mahasiswa
ADD CONSTRAINT fk_mahasiswa_dosen_pembimbing_id_dpa
    FOREIGN KEY (id_dpa)
    REFERENCES dosen(id_dosen)
    ON DELETE SET NULL;

ALTER TABLE mahasiswa
ADD CONSTRAINT fk_mahasiswa_kurikulum_kurikulum_id
    FOREIGN KEY (id_kurikulum)
    REFERENCES kurikulum(id_kurikulum)
    ON DELETE SET NULL;

ALTER TABLE mahasiswa
ADD CONSTRAINT fk_mahasiswa_user_id_user
    FOREIGN KEY (id_user)
    REFERENCES users(id_user)
    ON DELETE CASCADE;

ALTER TABLE kurikulum
ADD CONSTRAINT fk_kurikulum_program_studi_id_prodi
    FOREIGN KEY (id_prodi)
    REFERENCES program_studi(id_prodi)
    ON DELETE CASCADE;

ALTER TABLE detail_kurikulum
ADD CONSTRAINT fk_detail_kurikulum_kurikulum_id_kurikulum
    FOREIGN KEY (id_kurikulum)
    REFERENCES kurikulum(id_kurikulum)
    ON DELETE CASCADE;

ALTER TABLE detail_kurikulum
ADD CONSTRAINT fk_detail_kurikulum_matakuliah_id_matkul
    FOREIGN KEY (id_matkul)
    REFERENCES mata_kuliah(id_matkul)
    ON DELETE CASCADE;

ALTER TABLE kelas_ditawarkan
ADD CONSTRAINT fk_kelas_ditawarkan_periode_id_periode
    FOREIGN KEY (id_periode)
    REFERENCES periode_akademik(id_periode)
    ON DELETE CASCADE;

ALTER TABLE kelas_ditawarkan
ADD CONSTRAINT fk_kelas_ditawarkan_matakuliah_id_matkul
    FOREIGN KEY (id_matkul)
    REFERENCES mata_kuliah(id_matkul)
    ON DELETE CASCADE;

ALTER TABLE jadwal_kelas
ADD CONSTRAINT fk_jadwal_kelas_kelas_ditawarkan_id_kelas
    FOREIGN KEY (id_kelas)
    REFERENCES kelas_ditawarkan(id_kelas)
    ON DELETE CASCADE;

ALTER TABLE krs
ADD CONSTRAINT fk_krs_mahasiswa_id_mahasiswa
    FOREIGN KEY (id_mahasiswa)
    REFERENCES mahasiswa(id_mahasiswa)
    ON DELETE CASCADE;

ALTER TABLE krs
ADD CONSTRAINT fk_krs_periode_id_periode
    FOREIGN KEY (id_periode)
    REFERENCES periode_akademik(id_periode)
    ON DELETE CASCADE;
