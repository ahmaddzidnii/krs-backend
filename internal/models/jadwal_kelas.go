package models

import (
	"github.com/google/uuid"
	"time"
)

type JadwalKelas struct {
	IDJadwal     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IDKelas      uuid.UUID `gorm:"type:uuid;not null"`
	Hari         string    `gorm:"type:varchar(10);not null"`
	WaktuMulai   time.Time `gorm:"type:time;not null"`
	WaktuSelesai time.Time `gorm:"type:time;not null"`
	Ruang        string    `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	// --- DEFINISI RELASI ---
	// Relasi Many-to-One (sebuah jadwal dimiliki oleh satu kelas)
	// GORM akan menggunakan IDKelas sebagai foreign key.
	// ON DELETE CASCADE harus ditangani di level database atau melalui GORM hooks jika diperlukan.
	Kelas KelasDitawarkan `gorm:"foreignKey:IDKelas;references:IDKelas"`
}

func (JadwalKelas) TableName() string {
	return "jadwal_kelas"
}
