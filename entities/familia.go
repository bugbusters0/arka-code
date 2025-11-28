package entities

import "time"

// Familia representa el núcleo familiar en el sistema
// Estructura:
// - Correo: Identificador único de la familia (funciona como email principal)
// - Telefono: Número de contacto opcional de la familia (puede ser nulo)
// - Contraseña: Contraseña hash de acceso a la familia (no se serializa en JSON)
// - DeleteAt: Fecha de eliminación soft delete (nulo si no está eliminado)
// Uso: Agrupación principal de usuarios, gestión de acceso familiar
// Relaciones: Usuario (CorreoFamilia), Concepto (CorreoFamilia), Movimiento (CorreoFamilia)

type Familia struct {
	Correo     string     `json:"correo"`
	Telefono   *string    `json:"telefono"`
	Contraseña string     `json:"-"`
	DeleteAt   *time.Time `json:"delete_at"`
}
