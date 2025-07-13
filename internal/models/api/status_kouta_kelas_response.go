package api

type StatusKoutaKelasResponse struct {
	IDKelas         string `json:"id_kelas"`
	Terisi 	  int    `json:"terisi"`
	Kouta 	  int    `json:"kouta"`
	IsFull 	  bool   `json:"is_full"`
	IsJoined bool   `json:"is_joined"`
}
