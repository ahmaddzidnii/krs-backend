package domain

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// TimeOnly adalah tipe kustom untuk menangani tipe data TIME dari PostgreSQL.
type TimeOnly struct {
	sql.NullTime
}

// Scan mengimplementasikan interface sql.Scanner.
// Ini akan dipanggil oleh GORM saat membaca data dari database.
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		t.Time, t.Valid = time.Time{}, false
		return nil
	}

	// Driver pgx/pq bisa mengembalikan string atau time.Time
	switch v := value.(type) {
	case time.Time:
		t.Time, t.Valid = v, true
		return nil
	case []byte:
		// Coba parse dari format "15:04:05"
		parsedTime, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return err
		}
		t.Time, t.Valid = parsedTime, true
		return nil
	case string:
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		t.Time, t.Valid = parsedTime, true
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type TimeOnly", value)
}

// Value mengimplementasikan interface driver.Valuer.
func (t TimeOnly) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

type JadwalKelas struct {
	IDJadwal     uuid.UUID `gorm:"primary_key;column:id_jadwal" json:"id_jadwal"`
	IDKelas      uuid.UUID `gorm:"column:id_kelas;not null" json:"id_kelas"`
	Hari         string    `gorm:"column:hari;not null" json:"hari"`
	WaktuMulai   TimeOnly  `gorm:"column:waktu_mulai;not null" json:"waktu_mulai"`
	WaktuSelesai TimeOnly  `gorm:"column:waktu_selesai;not null" json:"waktu_selesai"`
	Ruang        string    `gorm:"column:ruang;not null" json:"ruang"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Kelas KelasDitawarkan `gorm:"foreignKey:id_kelas;references:id_kelas"`
}

func (JadwalKelas) TableName() string {
	return "jadwal_kelas"
}
