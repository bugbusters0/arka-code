package entities

import "time"

type Concepto struct {
	NombreConcepto string     `json:"nombre_concepto"`
	CorreoFamilia  string     `json:"correo_familia"`
	Tipo           int8       `json:"tipo"` // 0 = gasto, 1 = ingreso
	Icono          *string    `json:"icono"`
	Color          *string    `json:"color"`
	NombreUsuario  string     `json:"nombre_usuario"`
	DeleteAt       *time.Time `json:"delete_at"`
}

// Helpers
func (c *Concepto) TipoNombre() string {
	if c.Tipo == 1 {
		return "ingreso"
	}
	return "gasto"
}
