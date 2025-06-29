package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KRS merepresentasikan tabel 'krs'.
type KRS struct {
	IDKrs           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	IDMahasiswa     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_mahasiswa_periode"`
	IDPeriode       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_mahasiswa_periode"`
	TotalSksDiambil int       `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Mahasiswa Mahasiswa       `gorm:"foreignKey:IDMahasiswa"`
	Periode   PeriodeAkademik `gorm:"foreignKey:IDPeriode"`

	KelasDiambil []KelasDitawarkan `gorm:"many2many:detail_krs;foreignKey:IDKrs;joinForeignKey:IDKrs;References:IDKelas;joinReferences:IDKelas"`
}

func (KRS) TableName() string {
	return "krs"
}

type DetailKRS struct {
	IDKrs   uuid.UUID `gorm:"type:uuid;primaryKey"`
	IDKelas uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (DetailKRS) TableName() string {
	return "detail_krs"
}

// =================================================================================
// GORM HOOKS UNTUK MEREPLIKASI TRIGGER 'update_krs_summary'
// =================================================================================

// AfterCreate akan dipanggil setelah sebuah record DetailKRS berhasil dibuat.
func (dk *DetailKRS) AfterCreate(tx *gorm.DB) (err error) {
	return dk.updateKrsTotalSks(tx, "+")
}

// AfterDelete akan dipanggil setelah sebuah record DetailKRS berhasil dihapus.
func (dk *DetailKRS) AfterDelete(tx *gorm.DB) (err error) {
	return dk.updateKrsTotalSks(tx, "-")
}

// updateKrsTotalSks adalah fungsi helper yang berisi logika utama dari trigger.
func (dk *DetailKRS) updateKrsTotalSks(tx *gorm.DB, operator string) error {
	var sksToChange int

	// Langkah 1: Ambil SKS dari mata kuliah yang terkait dengan kelas yang diambil.
	// Ini mereplikasi: SELECT mk.sks INTO sks_to_change FROM kelas_ditawarkan ...
	err := tx.Model(&MataKuliah{}).
		Select("mata_kuliah.sks").
		Joins("JOIN kelas_ditawarkan ON kelas_ditawarkan.id_matkul = mata_kuliah.id_matkul").
		Where("kelas_ditawarkan.id_kelas = ?", dk.IDKelas).
		First(&sksToChange).Error

	if err != nil {
		// Jika kelas atau matkul tidak ditemukan, ini adalah error data.
		return fmt.Errorf("gagal mendapatkan sks untuk kelas %s: %w", dk.IDKelas, err)
	}

	if sksToChange == 0 {
		// Tidak ada perubahan SKS, tidak perlu update.
		return nil
	}

	// Langkah 2: Bangun ekspresi SQL untuk penambahan atau pengurangan.
	// gorm.Expr digunakan untuk membuat ekspresi SQL mentah, menghindari race condition.
	var updateExpr string
	if operator == "+" {
		updateExpr = "total_sks_diambil + ?"
	} else if operator == "-" {
		updateExpr = "total_sks_diambil - ?"
	} else {
		return fmt.Errorf("operator tidak valid: %s", operator)
	}

	// Langkah 3: Update kolom 'total_sks_diambil' di tabel KRS.
	// Ini mereplikasi: UPDATE krs SET total_sks_diambil = ...
	result := tx.Model(&KRS{}).
		Where("id_krs = ?", dk.IDKrs).
		Update("total_sks_diambil", gorm.Expr(updateExpr, sksToChange))

	if result.Error != nil {
		return fmt.Errorf("gagal mengupdate total_sks_diambil untuk krs %s: %w", dk.IDKrs, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("krs dengan id %s tidak ditemukan untuk diupdate", dk.IDKrs)
	}

	return nil
}
