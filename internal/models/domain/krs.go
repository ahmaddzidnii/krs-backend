package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KRS struct {
	IDKrs           uuid.UUID `gorm:"primaryKey;column:id_krs" json:"id_krs"`
	IDMahasiswa     uuid.UUID `gorm:"column:id_mahasiswa;not null;uniqueIndex:idx_mahasiswa_periode" json:"id_mahasiswa"`
	IDPeriode       uuid.UUID `gorm:"column:id_periode;not null;uniqueIndex:idx_mahasiswa_periode" json:"id_periode"`
	TotalSksDiambil int       `gorm:"column:total_sks_diambil;not null" json:"total_sks_diambil"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Mahasiswa Mahasiswa       `gorm:"foreignKey:IDMahasiswa" json:"mahasiswa"`
	Periode   PeriodeAkademik `gorm:"foreignKey:IDPeriode" json:"periode"`

	KelasDiambil []KelasDitawarkan `gorm:"many2many:detail_krs;foreignKey:id_krs;joinForeignKey:id_krs;References:id_kelas;joinReferences:id_kelas"`
}

func (KRS) TableName() string {
	return "krs"
}

type DetailKRS struct {
	IDKrs   uuid.UUID `gorm:"primaryKey;column:id_krs" json:"id_krs"`
	IDKelas uuid.UUID `gorm:"primaryKey;column:id_kelas" json:"id_kelas"`
}

func (DetailKRS) TableName() string {
	return "detail_krs"
}

// =================================================================================
// GORM HOOKS UNTUK MEREPLIKASI TRIGGER 'update_krs_summary'
// =================================================================================

func (dk *DetailKRS) AfterCreate(tx *gorm.DB) error {
	return dk.updateKrsTotalSks(tx, "+")
}

func (dk *DetailKRS) AfterDelete(tx *gorm.DB) error {
	return dk.updateKrsTotalSks(tx, "-")
}

func (dk *DetailKRS) updateKrsTotalSks(tx *gorm.DB, operator string) error {
	var sksToChange int

	err := tx.Model(&MataKuliah{}).
		Select("mata_kuliah.sks").
		Joins("JOIN kelas_ditawarkan ON kelas_ditawarkan.id_matkul = mata_kuliah.id_matkul").
		Where("kelas_ditawarkan.id_kelas = ?", dk.IDKelas).
		First(&sksToChange).Error

	if err != nil {
		return fmt.Errorf("gagal mendapatkan sks untuk kelas %s: %w", dk.IDKelas, err)
	}

	if sksToChange == 0 {
		return nil
	}

	var updateExpr string
	switch operator {
	case "+":
		updateExpr = "total_sks_diambil + ?"
	case "-":
		updateExpr = "total_sks_diambil - ?"
	default:
		return fmt.Errorf("operator tidak valid: %s", operator)
	}

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
