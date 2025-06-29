package repository

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DosenRepository interface {
	FindByUserID(userID uuid.UUID) (*models.Dosen, error)
}

type DosenRepositoryImpl struct {
	Db *gorm.DB
}

func NewDosenRepository(db *gorm.DB) DosenRepository {
	return &DosenRepositoryImpl{Db: db}
}

// FindByUserID sangat efisien karena mencari berdasarkan foreign key yang seharusnya di-index.
func (r *DosenRepositoryImpl) FindByUserID(userID uuid.UUID) (*models.Dosen, error) {
	var dosen models.Dosen
	// Di sini Anda bisa Preload data spesifik untuk mahasiswa jika perlu,
	// misalnya .Preload("ProgramStudi")
	err := r.Db.Where("id_user = ?", userID).First(&dosen).Error
	if err != nil {
		return nil, err
	}
	return &dosen, nil
}
