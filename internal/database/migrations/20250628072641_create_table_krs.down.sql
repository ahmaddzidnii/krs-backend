-- Migration Down Script
-- Perintah ini akan menghapus semua objek yang dibuat dalam urutan yang benar
-- untuk menghindari error karena adanya ketergantungan (dependency).
-- Urutan: Hapus Trigger -> Hapus Tabel -> Hapus Fungsi

-- 1. Hapus trigger terlebih dahulu karena bergantung pada tabel dan fungsi.
--    Trigger harus dihapus dari tabel tempat ia terpasang.
DROP TRIGGER IF EXISTS on_detail_krs_delete ON detail_krs;
DROP TRIGGER IF EXISTS on_detail_krs_insert ON detail_krs;
DROP TRIGGER IF EXISTS set_updated_at_krs ON krs;

-- 2. Hapus tabel. Tabel 'detail_krs' harus dihapus sebelum 'krs'
--    karena memiliki Foreign Key yang merujuk ke 'krs'.
DROP TABLE IF EXISTS detail_krs;
DROP TABLE IF EXISTS krs;

-- 3. Terakhir, hapus fungsi karena sudah tidak ada trigger yang menggunakannya.
DROP FUNCTION IF EXISTS update_krs_summary();

-- CATATAN: Script ini tidak menghapus fungsi 'trigger_set_updated_at()'.
-- Asumsinya, fungsi tersebut adalah fungsi umum yang mungkin digunakan oleh tabel lain.
-- Jika fungsi itu dibuat khusus hanya untuk tabel 'krs', Anda bisa menambahkannya di sini:
-- DROP FUNCTION IF EXISTS trigger_set_updated_at();