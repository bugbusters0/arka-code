package entities

import "time"

// Concepto representa una categoría para clasificar movimientos financieros
// Estructura:
// - NombreConcepto: Identificador único del concepto dentro de la familia
// - CorreoFamilia: Familia a la que pertenece el concepto (clave foránea)
// - Tipo: Clasificación del concepto (0 = gasto, 1 = ingreso)
// - Icono: Ícono representativo del concepto (opcional)
// - Color: Color identificador del concepto (opcional)
// - NombreUsuario: Usuario que creó el concepto
// - Activo: Estado de disponibilidad del concepto para nuevos movimientos
// - DeleteAt: Fecha de eliminación soft delete (nulo si no está eliminado)
// Uso: Categorización de movimientos, organización financiera
// Relaciones: Movimiento (NombreConcepto), PersonalizacionConcepto (NombreConcepto)

type Concepto struct {
	NombreConcepto string     `json:"nombre_concepto"`
	CorreoFamilia  string     `json:"correo_familia"`
	Tipo           int8       `json:"tipo"` // 0 = gasto, 1 = ingreso
	Icono          *string    `json:"icono"`
	Color          *string    `json:"color"`
	NombreUsuario  string     `json:"nombre_usuario"`
	Activo         bool       `json:"activo"`
	DeleteAt       *time.Time `json:"delete_at"`
}

// TipoNombre retorna el nombre legible del tipo de concepto
// Retorno: "gasto" para tipo 0, "ingreso" para tipo 1
// Uso: Mostrar el tipo en interfaces de usuario de forma legible
func (c *Concepto) TipoNombre() string {
	if c.Tipo == 1 {
		return "ingreso"
	}
	return "gasto"
}
