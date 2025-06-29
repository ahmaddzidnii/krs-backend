package models

import (
	"github.com/google/uuid"
	"time"
)

type Kurikulum struct {
	IDKurikulum   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id_kurikulum"`
	IDProdi       uuid.UUID `gorm:"type:uuid;not null" json:"id_prodi"`
	KodeKurikulum string    `gorm:"type:varchar(20);unique;not null" json:"kode_kurikulum"`
	Nama          string    `gorm:"type:varchar(100);not null" json:"nama"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	DetailKurikulum []DetailKurikulum `gorm:"foreignKey:id_kurikulum;references:id_kurikulum"`
}

func (m *Kurikulum) TableName() string {
	return "kurikulum"
}
