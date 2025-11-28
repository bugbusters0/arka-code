package entities

import "time"

// Movimiento representa un registro financiero (gasto o ingreso) en el sistema
// Estructura:
// - IdMovimiento: Identificador único autoincremental del movimiento
// - Fecha: Fecha en que se realizó el movimiento (obligatorio)
// - Monto: Valor monetario del movimiento (obligatorio, positivo)
// - Descripcion: Detalle opcional del movimiento (puede ser nulo)
// - NombreUsuario: Usuario que registró el movimiento (relación con tabla usuario)
// - NombreConcepto: Concepto categorizador del movimiento (relación con tabla concepto)
// - CorreoFamilia: Familia a la que pertenece el movimiento (relación con tabla familia)
// - DeleteAt: Fecha de eliminación soft delete (nulo si no está eliminado)
// Uso: Registro de transacciones financieras diarias, control de gastos e ingresos
// Relaciones: Usuario (NombreUsuario), Concepto (NombreConcepto), Familia (CorreoFamilia)
type Movimiento struct {
	IdMovimiento   int        `json:"id_movimiento"`
	Fecha          time.Time  `json:"fecha"`
	Monto          float64    `json:"monto"`
	Descripcion    *string    `json:"descripcion"`
	NombreUsuario  string     `json:"nombre_usuario"`
	NombreConcepto string     `json:"nombre_concepto"`
	CorreoFamilia  string     `json:"correo_familia"`
	DeleteAt       *time.Time `json:"delete_at"`
}
