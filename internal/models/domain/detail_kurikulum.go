package domain

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
	IDDetailKurikulum uuid.UUID       `gorm:"primaryKey;column:id_detail_kurikulum" json:"id_detail_kurikulum"`
	IDKurikulum       uuid.UUID       `gorm:"column:id_kurikulum" json:"id_kurikulum"`
	IDMatkul          uuid.UUID       `gorm:"column:id_matkul" json:"id_matkul"`
	JenisMatkul       JenisMataKuliah `gorm:"column:jenis_matkul" json:"jenis_matkul"`
	SemesterPaket     int             `gorm:"column:semester_paket" json:"semester_paket"`
	CreatedAt         time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Kurikulum Kurikulum  `gorm:"foreignKey:id_kurikulum;references:id_kurikulum" json:"kurikulum,omitempty"`
	Matkul    MataKuliah `gorm:"foreignKey:id_matkul;references:id_matkul" json:"matkul,omitempty"`
}

func (m *DetailKurikulum) TableName() string {
	return "detail_kurikulum"
}
