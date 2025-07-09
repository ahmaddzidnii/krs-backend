package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type MahasiswaRepository interface {
	FindByUserID(userID uuid.UUID) (*models.Mahasiswa, error)
	FindByNIM(nim string) (*models.Mahasiswa, error)
	FindByNIMWithTotalSKS(nim string) (*MahasiswaWithSKS, error)
}

type MahasiswaRepositoryImpl struct {
	Db          *gorm.DB
	RedisClient *redis.Client
	Context     context.Context
}

func NewMahasiswaRepository(db *gorm.DB, redisClient *redis.Client) MahasiswaRepository {
	return &MahasiswaRepositoryImpl{Db: db, RedisClient: redisClient, Context: context.Background()}
}

func (r *MahasiswaRepositoryImpl) FindByUserID(userID uuid.UUID) (*models.Mahasiswa, error) {
	cacheKey := "mahasiswa:userId:" + userID.String()
	cacheTTL := 1 * time.Hour

	// 1. Coba ambil dari Redis dulu
	cachedData, err := r.RedisClient.Get(r.Context, cacheKey).Result()
	if err == nil {
		// Cache HIT!
		var mahasiswa models.Mahasiswa
		if json.Unmarshal([]byte(cachedData), &mahasiswa) == nil {
			return &mahasiswa, nil
		}
	}
	// 2. Cache MISS! Ambil dari database.
	var mahasiswa models.Mahasiswa
	err = r.Db.Where("id_user = ?", userID).First(&mahasiswa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan
		}
		return nil, err // Error lain
	}
	// 3. Simpan hasil ke Redis untuk permintaan berikutnya.
	jsonData, err := json.Marshal(&mahasiswa)
	if err == nil {
		r.RedisClient.Set(r.Context, cacheKey, jsonData, cacheTTL).Err()
	}
	return &mahasiswa, nil
}

func (r *MahasiswaRepositoryImpl) FindByNIM(nim string) (*models.Mahasiswa, error) {
	cacheKey := "mahasiswa:nim:" + nim
	cacheTTL := 1 * time.Hour

	// 1. Coba ambil dari Redis dulu
	cachedData, err := r.RedisClient.Get(r.Context, cacheKey).Result()
	if err == nil {
		// Cache HIT!
		var mahasiswa models.Mahasiswa
		if json.Unmarshal([]byte(cachedData), &mahasiswa) == nil {
			return &mahasiswa, nil
		}
	}

	// 2. Cache MISS! Ambil dari database.
	var mahasiswa models.Mahasiswa
	err = r.Db.Where("nim = ?", nim).First(&mahasiswa).Error
	if err != nil {
		return nil, err // Termasuk jika record tidak ditemukan
	}

	// 3. Simpan hasil ke Redis untuk permintaan berikutnya.
	jsonData, err := json.Marshal(&mahasiswa)
	if err == nil {
		r.RedisClient.Set(r.Context, cacheKey, jsonData, cacheTTL).Err()
	}

	return &mahasiswa, nil
}

type MahasiswaWithSKS struct {
	models.Mahasiswa
	TotalSKSDiambil int `json:"total_sks_diambil"`
}

func (r *MahasiswaRepositoryImpl) FindByNIMWithTotalSKS(nim string) (*MahasiswaWithSKS, error) {
	var result MahasiswaWithSKS

	// Menggunakan COALESCE untuk memastikan SUM mengembalikan 0 jika hasilnya NULL
	err := r.Db.Model(&models.Mahasiswa{}).
		Select("mahasiswa.*, COALESCE(SUM(mata_kuliah.sks), 0) as total_sks_diambil").
		// Gunakan LEFT JOIN di sini
		Joins("LEFT JOIN krs ON krs.id_mahasiswa = mahasiswa.id_mahasiswa").
		Joins("LEFT JOIN detail_krs ON detail_krs.id_krs = krs.id_krs").
		Joins("LEFT JOIN kelas_ditawarkan ON kelas_ditawarkan.id_kelas = detail_krs.id_kelas").
		Joins("LEFT JOIN mata_kuliah ON mata_kuliah.id_matkul = kelas_ditawarkan.id_matkul").
		Where("mahasiswa.nim = ?", nim).
		Group("mahasiswa.id_mahasiswa").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}
