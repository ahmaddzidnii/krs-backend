package api

type DaftarPenawaranKelasResponse struct {
	IdKelas         string   `json:"id_kelas"`
	KodeMataKuliah  string   `json:"kode_mata_kuliah"`
	KodeKurikulum   string   `json:"kode_kurikulum"`
	NamaMataKuliah  string   `json:"nama_mata_kuliah"`
	JenisMataKuliah string   `json:"jenis_mata_kuliah"`
	Sks             int      `json:"sks"`
	SemesterPaket   int      `json:"semester_paket"`
	NamaKelas       string   `json:"nama_kelas"`
	DosenPengajar   []Dosen  `json:"dosen_pengajar"`
	Jadwal          []Jadwal `json:"jadwal"`
}

type PenawaranPerSemesterResponse struct {
	SemesterPaket map[string][]*DaftarPenawaranKelasResponse `json:"semester_paket"`
}

// Struct Dosen dan Jadwal (diasumsikan sudah ada)
type Dosen struct {
	NipDosen  string `json:"nip_dosen"`
	NamaDosen string `json:"nama_dosen"`
}

type Jadwal struct {
	Hari         string `json:"hari"`
	WaktuMulai   string `json:"waktu_mulai"`
	WaktuSelesai string `json:"waktu_selesai"`
	Ruangan      string `json:"ruangan"`
}
