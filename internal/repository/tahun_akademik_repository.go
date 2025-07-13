package repository

import (
	"context"
	"encoding/json"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/domain"
	"github.com/redis/go-redis/v9"
	"time"

	"gorm.io/gorm"
)

type TahunAkademikRepository interface {
	FindActive() (domain.PeriodeAkademik, error)
}

type TahunAkademikRepositoryImpl struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Context context.Context
}

func NewTahunAkademikRepository(db *gorm.DB, redisClient *redis.Client) TahunAkademikRepository {
	return &TahunAkademikRepositoryImpl{
		DB:      db,
		Redis:   redisClient,
		Context: context.Background(),
	}
}

func (r *TahunAkademikRepositoryImpl) FindActive() (domain.PeriodeAkademik, error) {
	cacheKey := "tahun_akademik:active"
	emptyResult := domain.PeriodeAkademik{}

	// 1. Coba ambil dari cache terlebih dahulu  caching
	result, err := r.Redis.Get(r.Context, cacheKey).Result()
	if err == nil {
		// Cache HIT: Data ditemukan
		var periodeAkademik domain.PeriodeAkademik
		if json.Unmarshal([]byte(result), &periodeAkademik) == nil {
			// Berhasil di-unmarshal, kembalikan data dari cache
			return periodeAkademik, nil
		}
	}

	// 2. Cache MISS: Ambil dari database
	var periodeAkademikFromDB domain.PeriodeAkademik
	err = r.DB.Where("is_active = ?", true).First(&periodeAkademikFromDB).Error
	if err != nil {
		return emptyResult, err
	}

	// 3. Simpan hasil dari database ke cache untuk api berikutnya
	jsonData, err := json.Marshal(periodeAkademikFromDB)
	if err == nil {
		// Set cache dengan waktu kedaluwarsa
		r.Redis.Set(r.Context, cacheKey, jsonData, 1*time.Hour)
	}

	return periodeAkademikFromDB, nil
}
