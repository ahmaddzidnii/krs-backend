package repository

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MahasiswaRepository interface {
	FindByUserID(userID uuid.UUID) (*models.Mahasiswa, error)
}

type MahasiswaRepositoryImpl struct {
	Db *gorm.DB
}

func NewMahasiswaRepository(db *gorm.DB) MahasiswaRepository {
	return &MahasiswaRepositoryImpl{Db: db}
}

// FindByUserID sangat efisien karena mencari berdasarkan foreign key yang seharusnya di-index.
func (r *MahasiswaRepositoryImpl) FindByUserID(userID uuid.UUID) (*models.Mahasiswa, error) {
	var mahasiswa models.Mahasiswa
	// Di sini Anda bisa Preload data spesifik untuk mahasiswa jika perlu,
	// misalnya .Preload("ProgramStudi")
	err := r.Db.Where("id_user = ?", userID).First(&mahasiswa).Error
	if err != nil {
		return nil, err
	}
	return &mahasiswa, nil
}
