package repository

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type AuthRepository interface {
	FindByCredential(credential string) (*domain.User, error)
}

type AuthRepositoryImpl struct {
	Db     *gorm.DB
	Logger *logrus.Logger
}

func NewAuthRepository(db *gorm.DB, logger *logrus.Logger) AuthRepository {
	return &AuthRepositoryImpl{
		Db:     db,
		Logger: logger,
	}
}

func (r *AuthRepositoryImpl) FindByCredential(credential string) (*domain.User, error) {
	log := r.Logger.WithField("credential", credential)
	log.Info("Mencari user berdasarkan kredential yang diberikan")
	var user domain.User
	startTime := time.Now()
	err := r.Db.Preload("Role").Where("username = ?", credential).First(&user).Error
	duration := time.Since(startTime)

	if duration > 5*time.Second {
		log.WithField("duration", duration).Warn("Pencarian  dengan credentials memakan waktu lebih dari 5 detik")
	} else {
		log.WithField("duration", duration).Info("Pencarian  selesai")
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}
