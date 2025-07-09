package models

import (
	"github.com/google/uuid"
	"time"
)

type JadwalKelas struct {
	IDJadwal     uuid.UUID `gorm:"primary_key;column:id_jadwal" json:"id_jadwal"`
	IDKelas      uuid.UUID `gorm:"column:id_kelas;not null" json:"id_kelas"`
	Hari         string    `gorm:"column:hari;not null" json:"hari"`
	WaktuMulai   time.Time `gorm:"column:waktu_mulai;not null" json:"waktu_mulai"`
	WaktuSelesai time.Time `gorm:"column:waktu_selesai;not null" json:"waktu_selesai"`
	Ruang        string    `gorm:"column:ruang;not null" json:"ruang"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Kelas KelasDitawarkan `gorm:"foreignKey:id_kelas;references:id_kelas"`
}

func (JadwalKelas) TableName() string {
	return "jadwal_kelas"
}
