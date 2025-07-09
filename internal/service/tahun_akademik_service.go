package service

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type TahunAkademikService interface {
	GetActiveTahunAkademik() (models.PeriodeAkademik, error)
}

type TahunAkademikServiceImpl struct {
	TahunAkademikRepository repository.TahunAkademikRepository
	Logger                  *logrus.Logger
}

func NewTahunAkademikService(repo repository.TahunAkademikRepository, logger *logrus.Logger) TahunAkademikService {
	return &TahunAkademikServiceImpl{
		TahunAkademikRepository: repo,
		Logger:                  logger,
	}
}

func (s *TahunAkademikServiceImpl) GetActiveTahunAkademik() (models.PeriodeAkademik, error) {
	return s.TahunAkademikRepository.FindActive()
}
