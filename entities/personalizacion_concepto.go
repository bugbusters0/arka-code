package entities

import "time"

// PersonalizacionConcepto representa la configuración personalizada de un concepto por usuario
// Estructura:
// - IdPersonalizacion: Identificador único autoincremental de la personalización
// - LimiteGasto: Límite máximo de gasto para el concepto (opcional, solo para gastos)
// - Activo: Estado de activación de la personalización
// - MontoPlanificado: Cantidad planificada para el concepto (opcional)
// - TipoPeriodoPlanificado: Frecuencia del monto planificado ("diario", "semanal", "mensual")
// - TipoPeriodoLimite: Frecuencia del límite de gasto ("diario", "semanal", "mensual")
// - DiaPeriodoPlanificado: Día específico para aplicar planificación (1-31, opcional)
// - DiaPeriodoLimite: Día específico para aplicar límite (1-31, opcional)
// - Notificacion: Habilitar notificaciones para este concepto
// - NombreUsuario: Usuario dueño de la personalización
// - NombreConcepto: Concepto personalizado
// - CorreoFamilia: Familia a la que pertenece
// - DeleteAt: Fecha de eliminación soft delete (nulo si no está eliminado)
// Uso: Configuración individual de conceptos, control presupuestario personalizado
// Relaciones: Concepto (NombreConcepto, CorreoFamilia), Usuario (NombreUsuario)

type PersonalizacionConcepto struct {
	IdPersonalizacion      int        `json:"id_personalizacion"`
	LimiteGasto            *float64   `json:"limite_gasto"`
	Activo                 bool       `json:"activo"`
	MontoPlanificado       *float64   `json:"monto_planificado"`
	TipoPeriodoPlanificado *string    `json:"tipo_periodo_planificado"` // "diario", "semanal", "mensual"
	TipoPeriodoLimite      *string    `json:"tipo_periodo_limite"`
	DiaPeriodoPlanificado  *int8      `json:"dia_periodo_planificado"` // 1-31
	DiaPeriodoLimite       *int8      `json:"dia_periodo_limite"`      // 1-31
	Notificacion           bool       `json:"notificacion"`
	NombreUsuario          string     `json:"nombre_usuario"`
	NombreConcepto         string     `json:"nombre_concepto"`
	CorreoFamilia          string     `json:"correo_familia"`
	DeleteAt               *time.Time `json:"delete_at"`
}
