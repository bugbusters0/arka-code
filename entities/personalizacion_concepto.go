package entities

import "time"

type PersonalizacionConcepto struct {
	IdPersonalizacion      int        `json:"id_personalizacion"`
	LimiteGasto            *float64   `json:"limite_gasto"`
	Activo                 bool       `json:"activo"`
	MontoPlanificado       *float64   `json:"monto_planificado"`
	TipoPeriodoPlanificado *string    `json:"tipo_periodo_planificado"` // "diario", "semanal", "mensual"
	TipoPeriodoLimite      *string    `json:"tipo_periodo_limite"`
	DiaPeriodoPlanificado  *int8      `json:"dia_periodo_planificado"` // 1-31
	Notificacion           bool       `json:"notificacion"`
	NombreUsuario          string     `json:"nombre_usuario"`
	NombreConcepto         string     `json:"nombre_concepto"`
	CorreoFamilia          string     `json:"correo_familia"`
	DeleteAt               *time.Time `json:"delete_at"`
}
