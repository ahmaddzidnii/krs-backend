package repository

import (
	"context"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models/domain"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PenjadwalanRepository interface {
	GetJadwalKelasDitawarkanBySemesterAndIdKurikulum(semesterPaket int, idKurikulum string) ([]*domain.KelasDitawarkan, error)
	GetKelasDitawarkanById(idKelas string) (*domain.KelasDitawarkan, error)
	GetIsJoinedKelasByNimAndIdKelas(idKelas string, nim string) (bool, error)
	GetJumlahTerisiKelasByIdKelas(idKelas string) (int64, error)

	GetKelasDitawarkanByIds(idKelas []string) ([]*domain.KelasDitawarkan, error)
	GetJumlahTerisiKelasByIds(idKelas []string) (map[string]int64, error)
	GetIsJoinedKelasForNim(idKelas []string, nim string) (map[string]bool, error)
}

type PenjadwalanRepositoryImpl struct {
	DB      *gorm.DB
	Redis   *redis.Client
	logger  *logrus.Logger
	Context context.Context
}

func NewPenjadwalanRepository(DB *gorm.DB, redis *redis.Client, logger *logrus.Logger) PenjadwalanRepository {
	return &PenjadwalanRepositoryImpl{DB: DB, Redis: redis, Context: context.Background(), logger: logger}
}

func (r *PenjadwalanRepositoryImpl) GetJadwalKelasDitawarkanBySemesterAndIdKurikulum(semesterPaket int, idKurikulum string) ([]*domain.KelasDitawarkan, error) {
	var result []*domain.KelasDitawarkan

	err := r.DB.Model(&domain.KelasDitawarkan{}).
		// Gunakan Joins HANYA untuk filtering
		Joins("JOIN public.periode_akademik pa ON pa.id_periode = kelas_ditawarkan.id_periode").
		Joins("JOIN public.mata_kuliah mk ON mk.id_matkul = kelas_ditawarkan.id_matkul").
		Joins("JOIN public.detail_kurikulum dk ON dk.id_matkul = mk.id_matkul").
		Joins("JOIN public.kurikulum kur ON kur.id_kurikulum = dk.id_kurikulum").

		// Terapkan kondisi WHERE dari tabel yang di-join
		Where("pa.is_active IS TRUE").
		//Where("dk.semester_paket = ?", semesterPaket).
		Where("kur.id_kurikulum = ?", idKurikulum).

		// Gunakan Preload untuk memuat data relasional (seperti 'include')
		// Ini akan menjalankan query terpisah untuk menghindari duplikasi
		Preload("JadwalKelas").
		Preload("DosenPengajar").
		// Preload bersarang (nested) untuk mendapatkan detail mata kuliah beserta kurikulumnya
		Preload("MataKuliah.DetailKurikulum.Kurikulum").

		// Ambil hasilnya
		Find(&result).Error

	if err != nil {
		r.logger.WithError(err).Error("Failed to get jadwal kelas ditawarkan by semester and id kurikulum")
		return nil, err
	}
	return result, nil
}

func (r *PenjadwalanRepositoryImpl) GetKelasDitawarkanById(idKelas string) (*domain.KelasDitawarkan, error) {
	var kelasDitawarkan domain.KelasDitawarkan

	err := r.DB.Model(&domain.KelasDitawarkan{}).
		Where("id_kelas = ?", idKelas).
		First(&kelasDitawarkan).Error

	if err != nil {
		r.logger.WithError(err).Error("Failed to get kouta kelas by id kelas")
		return nil, err
	}

	return &kelasDitawarkan, nil
}

func (r *PenjadwalanRepositoryImpl) GetIsJoinedKelasByNimAndIdKelas(idKelas string, nim string) (bool, error) {
	var isJoinedCount int64
	err := r.DB.Model(&domain.DetailKRS{}).
		Joins("JOIN krs ON detail_krs.id_krs = krs.id_krs").
		Joins("JOIN mahasiswa ON krs.id_mahasiswa = mahasiswa.id_mahasiswa").
		Where("detail_krs.id_kelas = ? AND mahasiswa.nim = ?", idKelas, nim).
		Count(&isJoinedCount).Error

	if err != nil {
		return false, err
	}

	isJoined := isJoinedCount > 0

	return isJoined, nil
}

func (r *PenjadwalanRepositoryImpl) GetJumlahTerisiKelasByIdKelas(idKelas string) (int64, error) {
	var jumlahTerisi int64
	err := r.DB.Model(&domain.DetailKRS{}).Where("id_kelas = ?", idKelas).Count(&jumlahTerisi).Error
	if err != nil {
		return 0, err
	}

	return jumlahTerisi, nil
}

// GetKelasDitawarkanByIds mengambil data beberapa kelas sekaligus menggunakan 'IN' clause.
func (r *PenjadwalanRepositoryImpl) GetKelasDitawarkanByIds(idKelas []string) ([]*domain.KelasDitawarkan, error) {
	var kelasDitawarkan []*domain.KelasDitawarkan
	if len(idKelas) == 0 {
		return kelasDitawarkan, nil
	}
	err := r.DB.Model(&domain.KelasDitawarkan{}).
		Where("id_kelas IN ?", idKelas).
		Find(&kelasDitawarkan).Error
	if err != nil {
		r.logger.WithError(err).Error("Failed to get kouta kelas by ids")
		return nil, err
	}
	return kelasDitawarkan, nil
}

type ResultJumlahTerisi struct {
	IDKelas string
	Jumlah  int64
}

// GetJumlahTerisiKelasByIds mengambil jumlah mahasiswa terisi untuk beberapa kelas sekaligus.
func (r *PenjadwalanRepositoryImpl) GetJumlahTerisiKelasByIds(idKelas []string) (map[string]int64, error) {
	jumlahTerisiMap := make(map[string]int64)
	if len(idKelas) == 0 {
		return jumlahTerisiMap, nil
	}

	var results []ResultJumlahTerisi
	err := r.DB.Model(&domain.DetailKRS{}).
		Select("id_kelas, count(*) as jumlah").
		Where("id_kelas IN ?", idKelas).
		Group("id_kelas").
		Scan(&results).Error

	if err != nil {
		r.logger.WithError(err).Error("Failed to get jumlah terisi in batch")
		return nil, err
	}

	for _, result := range results {
		jumlahTerisiMap[result.IDKelas] = result.Jumlah
	}

	return jumlahTerisiMap, nil
}

// GetIsJoinedKelasForNim mengecek status join mahasiswa pada beberapa kelas sekaligus.
func (r *PenjadwalanRepositoryImpl) GetIsJoinedKelasForNim(idKelas []string, nim string) (map[string]bool, error) {
	isJoinedMap := make(map[string]bool)
	if len(idKelas) == 0 {
		return isJoinedMap, nil
	}

	var joinedKelasIds []string
	err := r.DB.Model(&domain.DetailKRS{}).
		Joins("JOIN krs ON detail_krs.id_krs = krs.id_krs").
		Joins("JOIN mahasiswa ON krs.id_mahasiswa = mahasiswa.id_mahasiswa").
		Where("detail_krs.id_kelas IN ? AND mahasiswa.nim = ?", idKelas, nim).
		Pluck("detail_krs.id_kelas", &joinedKelasIds).Error

	if err != nil {
		return nil, err
	}

	for _, id := range joinedKelasIds {
		isJoinedMap[id] = true
	}

	return isJoinedMap, nil
}
