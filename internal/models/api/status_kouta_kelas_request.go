package api

type StatusKoutaKelasRequest struct {
	IDKelas string `json:"id_kelas" validate:"required,uuid"`
}
