package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type JenjangEnum string

const (
	Sarjana  JenjangEnum = "Sarjana (S1)"
	Magister JenjangEnum = "Magister (S2)"
)

func (j *JenjangEnum) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("gagal memindai StatusMahasiswa: nilai bukan string")
	}
	*j = JenjangEnum(str)
	return nil
}

func (j JenjangEnum) Value() (driver.Value, error) {
	return string(j), nil
}

type ProgramStudi struct {
	IDProdi    uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id_prodi"`
	IDFakultas uuid.UUID   `gorm:"type:uuid;not null" json:"id_fakultas"`
	KodeProdi  string      `gorm:"type:varchar(20);unique;not null" json:"kode_prodi"`
	Nama       string      `gorm:"type:varchar(100);not null" json:"nama"`
	Jenjang    JenjangEnum `gorm:"type:jenjang_enum;not null" json:"jenjang"`
	CreatedAt  time.Time   `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Fakultas Fakultas `gorm:"foreignKey:id_fakultas;references:id_fakultas"`

	Kurikulum []Kurikulum `gorm:"foreignKey:id_prodi;references:id_prodi"`
}

func (m *ProgramStudi) TableName() string {
	return "program_studi"
}
