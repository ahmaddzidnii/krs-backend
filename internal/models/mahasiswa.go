package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatusMahasiswa string

const (
	Aktif    StatusMahasiswa = "Aktif"
	Cuti     StatusMahasiswa = "Cuti"
	NonAktif StatusMahasiswa = "Non-Aktif"
)

func (s *StatusMahasiswa) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("gagal memindai StatusMahasiswa: nilai bukan string")
	}
	*s = StatusMahasiswa(str)
	return nil
}

func (s StatusMahasiswa) Value() (driver.Value, error) {
	return string(s), nil
}

type StatusPembayaran string

const (
	Lunas      StatusPembayaran = "Lunas"
	BelumLunas StatusPembayaran = "Belum Lunas"
)

func (s *StatusPembayaran) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("gagal memindai StatusPembayaran: nilai bukan string")
	}
	*s = StatusPembayaran(str)
	return nil
}

func (s StatusPembayaran) Value() (driver.Value, error) {
	return string(s), nil
}

type Mahasiswa struct {
	IDMahasiswa      uuid.UUID        `gorm:"primaryKey;column:id_mahasiswa" json:"id_mahasiswa"`
	IDUser           uuid.UUID        `gorm:"column:id_user;not null;unique" json:"id_user"`
	IDProdi          uuid.UUID        `gorm:"column:id_prodi;not null" json:"id_prodi"`
	IDDpa            uuid.UUID        `gorm:"column:id_dpa;not null" json:"id_dpa"`
	IDKurikulum      uuid.UUID        `gorm:"column:id_kurikulum;not null" json:"id_kurikulum"`
	NIM              string           `gorm:"column:nim;type:varchar(20);not null;unique" json:"nim"`
	Nama             string           `gorm:"column:nama;type:varchar(100);not null" json:"nama"`
	IPK              float64          `gorm:"column:ipk;type:numeric(3,2);not null;default:0.00" json:"ipk"`
	IPSLalu          float64          `gorm:"column:ips_lalu;type:numeric(3,2);not null;default:0.00" json:"ips_lalu"`
	SemesterBerjalan int              `gorm:"column:semester_berjalan;not null;default:1" json:"semester_berjalan"`
	SKSKumulatif     int              `gorm:"column:sks_kumulatif;not null;default:0" json:"sks_kumulatif"`
	JatahSKS         int              `gorm:"column:jatah_sks;not null;default:0" json:"jatah_sks"`
	StatusMahasiswa  StatusMahasiswa  `gorm:"column:status_mahasiswa;type:status_mahasiswa_enum;not null;default:'Aktif'" json:"status_mahasiswa"`
	StatusPembayaran StatusPembayaran `gorm:"column:status_pembayaran;type:status_pembayaran_enum;not null;default:'Belum Lunas'" json:"status_pembayaran"`
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:id_user;references:id_user" json:"user,omitempty"`
}

func (m *Mahasiswa) TableName() string {
	return "mahasiswa"
}

func (m *Mahasiswa) BeforeCreate(tx *gorm.DB) (err error) {
	m.hitungJatahSKS()
	if m.IDMahasiswa == uuid.Nil {
		m.IDMahasiswa = uuid.New()
	}
	return
}

func (m *Mahasiswa) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("IPSLalu") {
		m.hitungJatahSKS()
	}
	return
}

func (m *Mahasiswa) hitungJatahSKS() {
	switch {
	case m.IPSLalu >= 3.00:
		m.JatahSKS = 24
	case m.IPSLalu >= 2.50:
		m.JatahSKS = 22
	case m.IPSLalu >= 2.00:
		m.JatahSKS = 20
	case m.IPSLalu >= 1.50:
		m.JatahSKS = 18
	default:
		m.JatahSKS = 16
	}
}
