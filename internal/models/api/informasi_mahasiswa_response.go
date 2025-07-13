package api

type InformasiMahasiswaResponse struct {
	TahunAkademik string  `json:"tahun_akademik"`
	Semester      string  `json:"semester"`
	IPK           float64 `json:"ipk"`
	SksKumulatif  int     `json:"sks_kumulatif"`
	IpsLalu       float64 `json:"ips_lalu"`
	JatahSks      int     `json:"jatah_sks"`
	SksAmbil      int     `json:"sks_ambil"`
	SisaSks       int     `json:"sisa_sks"`
}
