package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type JenisSemesterEnum string

const (
	SemesterGanjil JenisSemesterEnum = "GANJIL"
	SemesterGenap  JenisSemesterEnum = "GENAP"
	SemesterPendek JenisSemesterEnum = "PENDEK"
)

func (j *JenisSemesterEnum) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("gagal memindai JenisSemesterEnum: nilai bukan string")
	}
	*j = JenisSemesterEnum(str)
	return nil
}

func (j JenisSemesterEnum) Value() (driver.Value, error) {
	return string(j), nil
}

type PeriodeAkademik struct {
	IDPeriode           uuid.UUID         `gorm:"primaryKey;column:id_periode" json:"id_periode"`
	TahunAkademik       string            `gorm:"column:tahun_akademik" json:"tahun_akademik"`
	JenisSemester       JenisSemesterEnum `gorm:"column:jenis_semester" json:"jenis_semester"`
	TanggalMulaiKRS     time.Time         `gorm:"column:tanggal_mulai_krs" json:"tanggal_mulai_krs"`
	TanggalSelesaiKRS   time.Time         `gorm:"column:tanggal_selesai_krs" json:"tanggal_selesai_krs"`
	JamMulaiHarianKRS   string            `gorm:"column:jam_mulai_harian_krs" json:"jam_mulai_harian_krs"`
	JamSelesaiHarianKRS string            `gorm:"column:jam_selesai_harian_krs" json:"jam_selesai_harian_krs"`
	IsActive            bool              `gorm:"column:is_active" json:"is_active"`
	CreatedAt           time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time         `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (m *PeriodeAkademik) TableName() string {
	return "periode_akademik"
}

func (j JenisSemesterEnum) String() string {
	switch j {
	case SemesterGanjil:
		return "GANJIL"
	case SemesterGenap:
		return "GENAP"
	case SemesterPendek:
		return "PENDEK"
	default:
		return "TIDAK DIKETAHUI"
	}
}
