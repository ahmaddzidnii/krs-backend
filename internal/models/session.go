package models

type RoleType struct {
	IDRole   string `json:"id_role"`
	RoleName string `json:"role_name"`
}
type Session struct {
	UserId     string   `json:"user_id"`
	NomorInduk string   `json:"nomor_induk"`
	Nama       string   `json:"nama"`
	Role       RoleType `json:"role"`
}
