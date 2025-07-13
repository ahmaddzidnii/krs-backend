package domain

import (
	"time"

	"github.com/google/uuid"
)

type Fakultas struct {
	IDFakultas   uuid.UUID `gorm:"primaryKey;column:id_fakultas" json:"id_fakultas"`
	KodeFakultas string    `gorm:"column:kode_fakultas" json:"kode_fakultas"`
	Nama         string    `gorm:"column:nama" json:"nama"`
	Singkatan    string    `gorm:"column:singkatan" json:"singkatan"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	ProgramStudi []ProgramStudi `gorm:"foreignKey:id_fakultas;references:id_fakultas"`
}

func (m *Fakultas) TableName() string {
	return "fakultas"
}
