package domain

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

func (s StatusPembayaran) String() string {
	switch s {
	case Lunas:
		return "Sudah Bayar"
	case BelumLunas:
		return "Belum Bayar"
	default:
		return "TIDAK DIKETAHUI"
	}
}

func (s StatusMahasiswa) String() string {
	switch s {
	case Aktif:
		return "Aktif"
	case Cuti:
		return "Cuti"
	case NonAktif:
		return "Non-Aktif"
	default:
		return "TIDAK DIKETAHUI"
	}
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
	NIM              string           `gorm:"column:nim;not null;unique" json:"nim"`
	Nama             string           `gorm:"column:nama;not null" json:"nama"`
	IPK              float64          `gorm:"column:ipk;not null" json:"ipk"`
	IPSLalu          float64          `gorm:"column:ips_lalu;not null" json:"ips_lalu"`
	SemesterBerjalan int              `gorm:"column:semester_berjalan;not null" json:"semester_berjalan"`
	SKSKumulatif     int              `gorm:"column:sks_kumulatif;not null" json:"sks_kumulatif"`
	JatahSKS         int              `gorm:"column:jatah_sks;not null" json:"jatah_sks"`
	StatusMahasiswa  StatusMahasiswa  `gorm:"column:status_mahasiswa;not null" json:"status_mahasiswa"`
	StatusPembayaran StatusPembayaran `gorm:"column:status_pembayaran;not null" json:"status_pembayaran"`
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:IDUser;references:IDUser" json:"user,omitempty"`
}

func (m *Mahasiswa) TableName() string {
	return "mahasiswa"
}

func (m *Mahasiswa) BeforeCreate(tx *gorm.DB) error {
	m.hitungJatahSKS()
	if m.IDMahasiswa == uuid.Nil {
		m.IDMahasiswa = uuid.New()
	}
	return nil
}

func (m *Mahasiswa) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("IPSLalu") {
		m.hitungJatahSKS()
	}
	return nil
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
