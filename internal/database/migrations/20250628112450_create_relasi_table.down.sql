-- Migration Down
-- Menghapus semua foreign key constraint yang telah ditambahkan.
-- Menggunakan "IF EXISTS" adalah praktik yang aman agar script bisa dijalankan
-- berulang kali tanpa menyebabkan error jika constraint sudah tidak ada.

ALTER TABLE program_studi
DROP CONSTRAINT IF EXISTS fk_fakultas_program_studi_id_fakultas;

ALTER TABLE mahasiswa
DROP CONSTRAINT IF EXISTS fk_mahasiswa_program_studi_id_prodi;

ALTER TABLE mahasiswa
DROP CONSTRAINT IF EXISTS fk_mahasiswa_dosen_pembimbing_id_dpa;

ALTER TABLE mahasiswa
DROP CONSTRAINT IF EXISTS fk_mahasiswa_kurikulum_kurikulum_id;

ALTER TABLE mahasiswa
DROP CONSTRAINT IF EXISTS fk_mahasiswa_user_id_user;

ALTER TABLE kurikulum
DROP CONSTRAINT IF EXISTS fk_kurikulum_program_studi_id_prodi;

ALTER TABLE detail_kurikulum
DROP CONSTRAINT IF EXISTS fk_detail_kurikulum_kurikulum_id_kurikulum;

ALTER TABLE detail_kurikulum
DROP CONSTRAINT IF EXISTS fk_detail_kurikulum_matakuliah_id_matkul;

ALTER TABLE kelas_ditawarkan
DROP CONSTRAINT IF EXISTS fk_kelas_ditawarkan_periode_id_periode;

ALTER TABLE kelas_ditawarkan
DROP CONSTRAINT IF EXISTS fk_kelas_ditawarkan_matakuliah_id_matkul;

ALTER TABLE jadwal_kelas
DROP CONSTRAINT IF EXISTS fk_jadwal_kelas_kelas_ditawarkan_id_kelas;

ALTER TABLE krs
DROP CONSTRAINT IF EXISTS fk_krs_mahasiswa_id_mahasiswa;

ALTER TABLE krs
DROP CONSTRAINT IF EXISTS fk_krs_periode_id_periode;