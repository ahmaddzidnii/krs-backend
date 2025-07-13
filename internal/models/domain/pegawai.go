package domain

import (
	"github.com/google/uuid"
	"time"
)

type Pegawai struct {
	IDPegawai uuid.UUID `gorm:"primaryKey;column:id_pegawai" json:"id_pegawai"`
	NIP       string    `gorm:"column:nip;unique;not null" json:"nip"`
	IDUser    uuid.UUID `gorm:"column:id_user;unique;not null" json:"id_user"`
	Nama      string    `gorm:"column:nama;not null" json:"nama"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:IDUser;references:IDUser" json:"user,omitempty"`
}

func (m *Pegawai) TableName() string {
	return "pegawai"
}
