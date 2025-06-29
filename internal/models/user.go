package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	IDUser    uuid.UUID `gorm:"primaryKey;column:id_user" json:"id_user"`
	IDRole    uuid.UUID `gorm:"column:id_role" json:"id_role"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Role      Role       `gorm:"foreignKey:id_role;references:id_role" json:"role"`
	Pegawai   *Pegawai   `gorm:"foreignKey:id_pegawai;references:id_user" json:"pegawai,omitempty"`
	Dosen     *Dosen     `gorm:"foreignKey:id_dosen;references:id_user" json:"dosen,omitempty"`
	Mahasiswa *Mahasiswa `gorm:"foreignKey:id_mahasiswa;references:id_user" json:"mahasiswa,omitempty"`
}

func (m *User) TableName() string {
	return "users"
}
