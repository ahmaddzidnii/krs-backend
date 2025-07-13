package handlers

import (
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/api"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/domain"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/service"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MahasiswaHandler struct {
	MahasiswaService   service.MahasiswaService
	PenjadwalanService service.PenjadwalanService
	Logger             *logrus.Logger
	Validator          *validator.Validate
}

func NewMahasiswaHandler(mahasiswaService service.MahasiswaService, penjadwalanService service.PenjadwalanService, logger *logrus.Logger, validator *validator.Validate) *MahasiswaHandler {
	return &MahasiswaHandler{
		MahasiswaService:   mahasiswaService,
		PenjadwalanService: penjadwalanService,
		Logger:             logger,
		Validator:          validator,
	}
}

func (m *MahasiswaHandler) GetSyaratPengisisanKRS(c *fiber.Ctx) error {
	sessionData, err := utils.GetLocals[domain.Session](c, "session_data")
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan data sesi")
		return utils.Error(c, fiber.StatusUnauthorized, "Gagal mendapatkan data sesi")
	}

	syaratPengisisanKrs, err := m.MahasiswaService.GetSyaratPengisianKRS(sessionData.NomorInduk)
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan syarat pengisian KRS")
		return utils.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan syarat pengisian KRS")
	}

	return utils.Success(c, fiber.StatusOK, "Berhasil mendapatkan syarat pengisian KRS", syaratPengisisanKrs)

}

func (m *MahasiswaHandler) GetInformasiUmum(c *fiber.Ctx) error {
	sessionData, err := utils.GetLocals[domain.Session](c, "session_data")
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan data sesi")
		return utils.Error(c, fiber.StatusUnauthorized, "Gagal mendapatkan data sesi")
	}
	informasiMahasiswa, err := m.MahasiswaService.GetInformasiMahasiswa(sessionData.NomorInduk)

	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan informasi mahasiswa")
		return utils.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan informasi mahasiswa")
	}

	return utils.Success(c, fiber.StatusOK, "Berhasil mendapatkan informasi mahasiswa", informasiMahasiswa)
}

func (m *MahasiswaHandler) GetPenawaranKelas(c *fiber.Ctx) error {
	sessionData, err := utils.GetLocals[domain.Session](c, "session_data")
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan data sesi")
		return utils.Error(c, fiber.StatusUnauthorized, "Gagal mendapatkan data sesi")
	}

	jadwalKelas, err := m.PenjadwalanService.GetPenawaranKelasByNim(sessionData.NomorInduk)
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan jadwal kelas")
		return utils.Error(c, fiber.StatusInternalServerError, "Gagal mendapatkan jadwal kelas")
	}

	return utils.Success(c, fiber.StatusOK, "Berhasil mendapatkan jadwal kelas", jadwalKelas)
}

func (m *MahasiswaHandler) GetStatusKoutaKelas(c *fiber.Ctx) error {
	var request api.StatusKoutaKelasRequest
	if err := c.BodyParser(&request); err != nil {
		m.Logger.WithError(err).Error("Gagal mem-parsing body requests login")
		return utils.Error(c, fiber.StatusBadRequest, "Invalid requests body")
	}

	if err := m.Validator.Struct(request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorBag := utils.GenerateValidationResponse(validationErrors)
			return utils.ValidationErrorResponse(c, errorBag)
		}
	}

	sessionData, err := utils.GetLocals[domain.Session](c, "session_data")
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan data sesi")
		return utils.Error(c, fiber.StatusUnauthorized, "Gagal mendapatkan data sesi")
	}

	statusKoutaKelas, err := m.PenjadwalanService.GetStatusKoutaKelasByIdKelas(request.IDKelas, &sessionData)

	if err != nil {
		return err
	}

	return utils.Success(c, fiber.StatusOK, "Berhasil mendapatkan status kouta kelas", statusKoutaKelas)
}

func (m *MahasiswaHandler) GetStatusKoutaKelasBatch(c *fiber.Ctx) error {
	type StatusKoutaKelasBatchRequest struct {
		IDKelas []string `json:"id_kelas" validate:"required,dive,required"`
	}

	var request *StatusKoutaKelasBatchRequest

	if err := c.BodyParser(&request); err != nil {
		m.Logger.WithError(err).Error("Gagal mem-parsing body request batch")
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := m.Validator.Struct(request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorBag := utils.GenerateValidationResponse(validationErrors)
			return utils.ValidationErrorResponse(c, errorBag)
		}
	}

	sessionData, err := utils.GetLocals[domain.Session](c, "session_data")
	if err != nil {
		m.Logger.WithError(err).Error("Gagal mendapatkan data sesi")
		return utils.Error(c, fiber.StatusUnauthorized, "Gagal mendapatkan data sesi")
	}

	// Panggil service method yang baru
	statusKoutaKelas, err := m.PenjadwalanService.GetStatusKoutaKelasInBatch(request.IDKelas, &sessionData)
	if err != nil {
		// Asumsikan service akan mengembalikan error yang sesuai (misal: internal server error)
		return err
	}

	return utils.Success(c, fiber.StatusOK, "Berhasil mendapatkan status kouta kelas secara batch", statusKoutaKelas)
}
