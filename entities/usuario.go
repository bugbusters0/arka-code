package entities

import "time"

// Usuario representa un miembro individual de una familia en el sistema
// Estructura:
// - NombreUsuario: Identificador único del usuario dentro del sistema
// - Rol: Nivel de privilegios (0 = miembro, 1 = administrador)
// - ContrasenaPersonal: Contraseña hash personal del usuario (no se serializa en JSON)
// - NombrePersonal: Nombre real del usuario para mostrar
// - CorreoFamilia: Familia a la que pertenece el usuario (clave foránea)
// - DeleteAt: Fecha de eliminación soft delete (nulo si no está eliminado)
// Uso: Gestión de usuarios, control de acceso individual, personalización
// Relaciones: Familia (CorreoFamilia), Movimiento (NombreUsuario), Concepto (NombreUsuario)
type Usuario struct {
	NombreUsuario      string     `json:"nombre_usuario"`
	Rol                int8       `json:"rol"` // 0 = member, 1 = admin
	ContrasenaPersonal string     `json:"-"`
	NombrePersonal     string     `json:"nombre_personal"`
	CorreoFamilia      string     `json:"correo_familia"`
	DeleteAt           *time.Time `json:"delete_at"`
}

// IsAdmin verifica si el usuario tiene privilegios de administrador
// Retorno: true si el rol es 1 (administrador), false en caso contrario
// Uso: Control de acceso a nivel de entidad, validaciones de negocio
func (u *Usuario) IsAdmin() bool {
	return u.Rol == 1
}

// RoleName retorna el nombre legible del rol del usuario
// Retorno: "admin" para rol 1, "member" para rol 0
// Uso: Mostrar el rol en interfaces de usuario, logs y reportes
func (u *Usuario) RoleName() string {
	if u.Rol == 1 {
		return "admin"
	}
	return "member"
}
