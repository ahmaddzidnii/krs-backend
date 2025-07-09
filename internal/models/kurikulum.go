package models

import (
	"github.com/google/uuid"
	"time"
)

type Kurikulum struct {
	IDKurikulum   uuid.UUID `gorm:"primaryKey;column:id_kurikulum" json:"id_kurikulum"`
	IDProdi       uuid.UUID `gorm:"column:id_prodi;not null" json:"id_prodi"`
	KodeKurikulum string    `gorm:"column:kode_kurikulum;unique;not null" json:"kode_kurikulum"`
	Nama          string    `gorm:"column:nama;not null" json:"nama"`
	IsActive      bool      `gorm:"column:is_active;not null" json:"is_active"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	ProgramStudi    ProgramStudi      `gorm:"foreignKey:id_prodi;references:id_prodi" json:"program_studi,omitempty"`
	DetailKurikulum []DetailKurikulum `gorm:"foreignKey:id_kurikulum;references:id_kurikulum" json:"detail_kurikulum"`
}

func (m *Kurikulum) TableName() string {
	return "kurikulum"
}
