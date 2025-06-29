package repository

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PegawaiRepository interface {
	FindByUserID(userID uuid.UUID) (*models.Pegawai, error)
}

type PegawaiRepositoryImpl struct {
	Db *gorm.DB
}

func NewPegawaiRepository(db *gorm.DB) PegawaiRepository {
	return &PegawaiRepositoryImpl{Db: db}
}

func (r *PegawaiRepositoryImpl) FindByUserID(userID uuid.UUID) (*models.Pegawai, error) {
	var pegawai models.Pegawai
	err := r.Db.Where("id_user = ?", userID).First(&pegawai).Error
	if err != nil {
		return nil, err
	}
	return &pegawai, nil
}
