package models

import (
	"github.com/google/uuid"
	"time"
)

type Dosen struct {
	IDDosen   uuid.UUID `gorm:"primaryKey;column:id_dosen" json:"id_dosen"`
	NIP       string    `gorm:"column:nip;type:varchar(18);unique;not null" json:"nip"`
	IDUser    uuid.UUID `gorm:"column:id_user;unique;not null" json:"id_user"`
	Nama      string    `gorm:"column:nama;type:varchar(100);not null" json:"nama"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	User        User              `gorm:"foreignKey:id_user;references:id_user" json:"user,omitempty"`
	KelasDiampu []KelasDitawarkan `gorm:"many2many:dosen_pengajar_kelas;foreignKey:id_dosen;joinForeignKey:id_dosen;References:id_kelas;joinReferences:id_kelas"`
}

func (m *Dosen) TableName() string {
	return "dosen"
}
