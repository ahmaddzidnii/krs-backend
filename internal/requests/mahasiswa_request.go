package requests

type GetInformasiUmumRequest struct {
	Nim string `json:"nim" validate:"required"`
}
