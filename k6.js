import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Trend } from 'k6/metrics';

// --- Konfigurasi Uji Beban ---
export const options = {
    stages: [
        { duration: '10s', target: 1 },  // Naik perlahan ke 100 pengguna virtual selama 30 detik
        { duration: '1m', target: 1000 },   // Bertahan di 100 pengguna selama 1 menit (beban puncak)
        { duration: '30s', target: 0 },    // Turun kembali ke 0 pengguna selama 30 detik
    ],
    thresholds: {
        // 95% dari semua request harus selesai di bawah 800ms
        http_req_duration: ['p(95)<800'],
        // Kurang dari 1% dari request yang boleh gagal
        http_req_failed: ['rate<0.01'],
        // Tentukan kriteria keberhasilan untuk setiap proses
        'login_duration': ['p(95)<800'],
        'session_check_duration': ['p(95)<500'], // Metrik baru untuk cek sesi
        'logout_duration': ['p(95)<500'],
    },
};

// --- Metrik Kustom ---
const loginDuration = new Trend('login_duration');
const sessionCheckDuration = new Trend('session_check_duration'); // Metrik baru
const logoutDuration = new Trend('logout_duration');

// --- Data Uji ---
// Dalam skenario nyata, data ini sebaiknya dibaca dari file CSV.
const testUsers = [
    { nim: '23106050077', password: '12345678' },
    // { nim: '87654321', password: 'password456' },
    // Tambahkan lebih banyak user untuk simulasi yang lebih baik
];

// --- Skenario Uji Beban (Lifecycle per Virtual User) ---
export default function () {
    // Ambil base URL dari environment variable, dengan fallback.
    const BASE_URL = __ENV.BASE_URL || 'https://apigo.masako.my.id';

    // Setiap Virtual User (VU) akan memilih satu akun dari data uji.
    const user = testUsers[__VU % testUsers.length];

    // Variabel untuk menyimpan session ID
    let sessionId = '';

    // Gunakan group() untuk mengelompokkan request dalam laporan hasil.
    group('1. Proses Login', function () {
        const loginPayload = JSON.stringify({
            username: user.nim, // Payload disesuaikan dengan flow Anda
            password: user.password,
        });

        const params = {
            headers: { 'Content-Type': 'application/json' },
        };

        const res = http.post(`${BASE_URL}/api/v1/auth/login`, loginPayload, params);
        loginDuration.add(res.timings.duration);


        const loginSuccess = check(res, {
            'Login berhasil (status 200)': (r) => r.status === 200,
            'Respons berisi session_id': (r) => r.json('data.session_id') !== undefined && r.json('data.session_id') !== null,
        });

        // Jika login berhasil, simpan session ID untuk request selanjutnya.
        if (loginSuccess) {
            sessionId = res.json('data.session_id');
        }
    });

    // Lanjutkan hanya jika login berhasil dan mendapatkan session ID
    if (sessionId) {
        // Jeda untuk mensimulasikan pengguna melihat-lihat halaman setelah login
        sleep(Math.random() * 2 + 1); // Jeda acak antara 1 hingga 3 detik

        const authHeaders = {
            headers: {
                // Sesuaikan format header otentikasi jika berbeda,
                // contoh: 'X-Session-ID': sessionId
                'Authorization': `Bearer ${sessionId}`,
            },
        };

        group('2. Cek Sesi Pengguna', function () {
            // Kirim request GET ke endpoint session
            const res = http.get(`${BASE_URL}/api/v1/auth/session`, authHeaders);
            sessionCheckDuration.add(res.timings.duration);

            check(res, {
                'Cek Sesi berhasil (status 200)': (r) => r.status === 200,
                'Data NIM sesuai': (r) => r.json('data.nim') === user.nim,
            });
        });

        // Jeda sebelum logout
        sleep(1);

        group('3. Proses Logout', function () {
            // Menggunakan header yang sama dengan Cek Sesi
            // Diasumsikan endpoint logout memerlukan otentikasi
            const res = http.post(`${BASE_URL}/api/v1/auth/logout`, null, authHeaders);
            logoutDuration.add(res.timings.duration);

            check(res, {
                'Logout berhasil (status 200)': (r) => r.status === 200,
                'Pesan logout sesuai': (r) => r.json('data.message') === 'Logout successful',
            });
        });

    } else {
        // Jika login gagal, jangan lanjutkan.
        console.error(`Login gagal untuk NIM: ${user.nim}. VU berhenti.`);
    }

    // Beri jeda singkat sebelum iterasi berikutnya oleh VU yang sama
    sleep(1);
}