package domain

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
		return fmt.Errorf("gagal memindai JenjangEnum: nilai bukan string")
	}
	*j = JenjangEnum(str)
	return nil
}

func (j JenjangEnum) Value() (driver.Value, error) {
	return string(j), nil
}

type ProgramStudi struct {
	IDProdi    uuid.UUID   `gorm:"primaryKey;column:id_prodi" json:"id_prodi"`
	IDFakultas uuid.UUID   `gorm:"column:id_fakultas" json:"id_fakultas"`
	KodeProdi  string      `gorm:"column:kode_prodi;unique" json:"kode_prodi"`
	Nama       string      `gorm:"column:nama" json:"nama"`
	Jenjang    JenjangEnum `gorm:"column:jenjang" json:"jenjang"`
	CreatedAt  time.Time   `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Fakultas  Fakultas    `gorm:"foreignKey:id_fakultas;references:id_fakultas" json:"fakultas,omitempty"`
	Kurikulum []Kurikulum `gorm:"foreignKey:id_prodi;references:id_prodi" json:"kurikulum,omitempty"`
}

func (m *ProgramStudi) TableName() string {
	return "program_studi"
}
