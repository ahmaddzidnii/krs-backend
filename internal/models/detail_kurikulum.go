package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type JenisMataKuliah string

const (
	Wajib   JenisMataKuliah = "Wajib"
	Pilihan JenisMataKuliah = "Pilihan"
)

func (j *JenisMataKuliah) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("gagal memindai StatusPembayaran: nilai bukan string")
	}
	*j = JenisMataKuliah(str)
	return nil
}

func (j JenisMataKuliah) Value() (driver.Value, error) {
	return string(j), nil
}

type DetailKurikulum struct {
	IDDetailKurikulum uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id_detail_kurikulum"`
	IDKurikulum       uuid.UUID       `gorm:"type:uuid;not null" json:"id_kurikulum"`
	IDMatkul          uuid.UUID       `gorm:"type:uuid;not null" json:"id_matkul"`
	JenisMatkul       JenisMataKuliah `gorm:"type:jenis_mata_kuliah;not null" json:"jenis_matkul"`
	SemesterPaket     int             `gorm:"not null" json:"semester_paket"`
	CreatedAt         time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (m *DetailKurikulum) TableName() string {
	return "detail_kurikulum"
}
