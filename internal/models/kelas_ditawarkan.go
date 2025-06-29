package models

import (
	"github.com/google/uuid"
	"time"
)

type KelasDitawarkan struct {
	IDKelas   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id_kelas"`
	IDPeriode uuid.UUID `gorm:"type:uuid;not null" json:"id_periode"`
	IDMatkul  uuid.UUID `gorm:"type:uuid;not null" json:"id_matkul"`
	Kouta     int       `gorm:"not null" json:"kouta"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	PeriodeAkademik PeriodeAkademik `gorm:"foreignKey:IDPeriode" json:"periode_akademik"`
	MataKuliah      MataKuliah      `gorm:"foreignKey:IDMatkul" json:"mata_kuliah"`
	DosenPengajar   []Dosen         `gorm:"many2many:dosen_pengajar_kelas;foreignKey:IDKelas;joinForeignKey:IDKelas;References:IDDosen;joinReferences:IDDosen"`
}

func (m *KelasDitawarkan) TableName() string {
	return "kelas_ditawarkan"
}
