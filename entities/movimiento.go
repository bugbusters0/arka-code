package entities

import "time"

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
