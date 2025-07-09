package handlers

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/service"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MahasiswaHandler struct {
	MahasiswaService service.MahasiswaService
	Logger           *logrus.Logger
	Validator        *validator.Validate
}

func NewMahasiswaHandler(mahasiswaService service.MahasiswaService, logger *logrus.Logger, validator *validator.Validate) *MahasiswaHandler {
	return &MahasiswaHandler{
		MahasiswaService: mahasiswaService,
		Logger:           logger,
		Validator:        validator,
	}
}

func (m *MahasiswaHandler) GetSyaratPengisisanKRS(c *fiber.Ctx) error {
	sessionData, err := utils.GetLocals[models.Session](c, "session_data")
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
	sessionData, err := utils.GetLocals[models.Session](c, "session_data")
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
