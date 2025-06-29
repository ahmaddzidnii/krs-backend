package models

import (
	"github.com/google/uuid"
	"time"
)

type MataKuliah struct {
	IDMatkul   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id_matkul"`
	KodeMatkul string    `gorm:"type:varchar(10);unique;not null" json:"kode_matkul"`
	Nama       string    `gorm:"type:varchar(100);not null" json:"nama"`
	SKS        int       `gorm:"not null;check:sks > 0" json:"sks"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (m *MataKuliah) TableName() string {
	return "mata_kuliah"
}
