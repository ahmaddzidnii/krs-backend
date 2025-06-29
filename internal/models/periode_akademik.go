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
		return fmt.Errorf("gagal memindai StatusMahasiswa: nilai bukan string")
	}
	*j = JenisSemesterEnum(str)
	return nil
}

func (j JenisSemesterEnum) Value() (driver.Value, error) {
	return string(j), nil
}

type PeriodeAkademik struct {
	IDPeriode           uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id_periode"`
	TahunAkademik       string            `gorm:"type:varchar(9);unique;not null" json:"tahun_akademik"`
	JenisSemester       JenisSemesterEnum `gorm:"type:jenis_semester_enum;not null" json:"jenis_semester"`
	TanggalMulaiKRS     time.Time         `gorm:"type:date;not null" json:"tanggal_mulai_krs"`
	TanggalSelesaiKRS   time.Time         `gorm:"type:date;not null" json:"tanggal_selesai_krs"`
	JamMulaiHarianKRS   string            `gorm:"type:time;not null" json:"jam_mulai_harian_krs"`   // bisa juga pakai time.Time
	JamSelesaiHarianKRS string            `gorm:"type:time;not null" json:"jam_selesai_harian_krs"` // bisa juga pakai time.Time
	IsActive            bool              `gorm:"not null;default:false" json:"is_active"`
	CreatedAt           time.Time         `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time         `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (m *PeriodeAkademik) TableName() string {
	return "periode_akademik"
}
