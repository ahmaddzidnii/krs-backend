package domain

import (
	"github.com/google/uuid"
	"time"
)

type MataKuliah struct {
	IDMatkul   uuid.UUID `gorm:"primaryKey;column:id_matkul" json:"id_matkul"`
	KodeMatkul string    `gorm:"column:kode_matkul;unique;not null" json:"kode_matkul"`
	Nama       string    `gorm:"column:nama;not null" json:"nama"`
	SKS        int       `gorm:"column:sks;not null" json:"sks"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	DetailKurikulum []DetailKurikulum `gorm:"foreignKey:id_matkul;references:id_matkul" json:"detail_kurikulum,omitempty"`
}

func (m *MataKuliah) TableName() string {
	return "mata_kuliah"
}
