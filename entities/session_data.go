package entities

// SessionData almacena la información de la sesión del usuario activo
// Estructura:
// - NombreUsuario: Identificador único del usuario en sesión
// - CorreoFamilia: Familia a la que pertenece el usuario
// - Rol: Privilegios del usuario (0 = miembro, 1 = administrador)
// - NombrePersonal: Nombre real del usuario para mostrar
// Uso: Gestión de sesiones, control de acceso, personalización de interfaces
// Métodos: Proporciona helpers para verificar roles y mostrar información
type SessionData struct {
	NombreUsuario  string
	CorreoFamilia  string
	Rol            int8
	NombrePersonal string
}

// IsAdmin verifica si el usuario en sesión tiene privilegios de administrador
// Retorno: true si el rol es 1 (administrador), false en caso contrario
// Uso: Control de acceso a funcionalidades administrativas
func (s *SessionData) IsAdmin() bool {
	return s.Rol == 1
}

// RoleName retorna el nombre legible del rol del usuario
// Retorno: "admin" para rol 1, "member" para rol 0
// Uso: Mostrar el rol en interfaces de usuario de forma legible
func (s *SessionData) RoleName() string {
	if s.Rol == 1 {
		return "admin"
	}
	return "member"
}
