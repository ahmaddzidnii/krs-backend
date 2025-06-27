-- Migrasi 'down' ini akan menghapus semua objek yang dibuat oleh file 'up'
-- dalam urutan terbalik untuk memastikan tidak ada error dependensi.

-- STEP 1: Hapus trigger dari tabel 'fakultas'.
-- Trigger bergantung pada tabel dan fungsi, jadi ini harus dihapus terlebih dahulu.
DROP TRIGGER IF EXISTS set_updated_at_fakultas ON fakultas;

-- STEP 2: Hapus tabel 'fakultas'.
-- Tabel harus dihapus sebelum fungsi yang mungkin bergantung padanya.
DROP TABLE IF EXISTS fakultas;