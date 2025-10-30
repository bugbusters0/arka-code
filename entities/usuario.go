package entities

import "time"

type Usuario struct {
	NombreUsuario      string     `json:"nombre_usuario"`
	Rol                int8       `json:"rol"` // 0 = member, 1 = admin
	ContrasenaPersonal string     `json:"-"`
	NombrePersonal     string     `json:"nombre_personal"`
	CorreoFamilia      string     `json:"correo_familia"`
	DeleteAt           *time.Time `json:"delete_at"`
}

func (u *Usuario) IsAdmin() bool {
	return u.Rol == 1
}

func (u *Usuario) RoleName() string {
	if u.Rol == 1 {
		return "admin"
	}
	return "member"
}
