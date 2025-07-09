package handlers

// Handlers struct menampung semua instance handler aplikasi.
type Handlers struct {
	AuthHandler      *AuthHandler
	MahasiswaHandler *MahasiswaHandler
	// Tambahkan handler baru di sini di masa depan
	// DosenHandler *DosenHandler
}

// NewHandlers adalah provider untuk struct Handlers.
// Fungsi ini akan dipanggil oleh Google Wire untuk membuat instance Handlers.
func NewHandlers(auth *AuthHandler, mhs *MahasiswaHandler) *Handlers {
	return &Handlers{
		AuthHandler:      auth,
		MahasiswaHandler: mhs,
	}
}
