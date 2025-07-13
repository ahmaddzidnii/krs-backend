package api

type SyaratPengisisanKrsResponse struct {
	Judul                string       `json:"judul"`
	DataSyarat           []SyaratItem `json:"data_syarat"`
	PengisisanKrsEnabled bool         `json:"pengisisan_krs_enabled"`
}

type SyaratItem struct {
	Syarat string `json:"syarat"`
	Isi    string `json:"isi"`
	Status bool   `json:"status"`
}
