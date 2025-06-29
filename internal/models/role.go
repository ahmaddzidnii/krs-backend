package models

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	IDRole    uuid.UUID `gorm:"primaryKey;column:id_role" json:"id_role"`
	RoleName  string    `gorm:"column:role_name;unique;not null" json:"role_name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Users []User `gorm:"foreignKey:id_role;references:id_role" json:"users"`
}

func (m *Role) TableName() string {
	return "roles"
}
