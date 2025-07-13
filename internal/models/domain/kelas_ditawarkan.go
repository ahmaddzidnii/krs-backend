package domain

import (
	"github.com/google/uuid"
	"time"
)

type KelasDitawarkan struct {
	IDKelas   uuid.UUID `gorm:"primaryKey;column:id_kelas" json:"id_kelas"`
	IDPeriode uuid.UUID `gorm:"column:id_periode;not null" json:"id_periode"`
	IDMatkul  uuid.UUID `gorm:"column:id_matkul;not null" json:"id_matkul"`
	NamaKelas string    `gorm:"column:nama_kelas;not null" json:"nama_kelas"`
	Kouta     int       `gorm:"column:kouta;not null" json:"kouta"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	PeriodeAkademik PeriodeAkademik `gorm:"foreignKey:id_periode;references:id_periode" json:"periode_akademik"`

	MataKuliah MataKuliah `gorm:"foreignKey:id_matkul;references:id_matkul" json:"mata_kuliah"`

	JadwalKelas []JadwalKelas `gorm:"foreignKey:id_kelas;references:id_kelas" json:"jadwal_kelas,omitempty"`

	DosenPengajar []Dosen `gorm:"many2many:dosen_pengajar_kelas;foreignKey:id_kelas;joinForeignKey:id_kelas;References:id_dosen;joinReferences:id_dosen"`

	KRS []KRS `gorm:"many2many:detail_krs;foreignKey:id_kelas;joinForeignKey:id_kelas;References:id_krs;joinReferences:id_krs"`
}

func (m *KelasDitawarkan) TableName() string {
	return "kelas_ditawarkan"
}
