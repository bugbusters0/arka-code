package entities

// Estructura para datos de sesión
type SessionData struct {
	NombreUsuario  string
	CorreoFamilia  string
	Rol            int8
	NombrePersonal string
}

func (s *SessionData) IsAdmin() bool {
	return s.Rol == 1
}

func (s *SessionData) RoleName() string {
	if s.Rol == 1 {
		return "admin"
	}
	return "member"
}
