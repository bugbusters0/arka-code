package validators

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MovimientoValidator struct {
	*ValidatorBase
}

var MovimientoValidatorInstance = &MovimientoValidator{}

func (v *MovimientoValidator) Validate(r *http.Request) ValidationResult {
	result := NewValidationResult()

	// Validar monto
	montoStr := strings.TrimSpace(r.FormValue("monto"))
	if montoStr == "" {
		result.Errors["monto"] = "El monto es obligatorio"
		result.Success = false
	} else {
		monto, err := strconv.ParseFloat(montoStr, 64)
		if err != nil || monto <= 0 {
			result.Errors["monto"] = "El monto debe ser un número mayor a 0"
			result.Success = false
		} else {
			result.CleanData["monto"] = monto
		}
	}

	// Validar concepto
	concepto := strings.TrimSpace(r.FormValue("nombreConcepto"))
	if concepto == "" {
		result.Errors["nombreConcepto"] = "Debe seleccionar un concepto"
		result.Success = false
	} else {
		result.CleanData["nombreConcepto"] = concepto
	}

	// Validar descripción (opcional, máximo 500 caracteres)
	descripcion := strings.TrimSpace(r.FormValue("descripcion"))
	if len(descripcion) > 500 {
		result.Errors["descripcion"] = "La descripción no puede exceder 500 caracteres"
		result.Success = false
	} else if descripcion != "" {
		result.CleanData["descripcion"] = descripcion
	}

	// Validar fecha
	fechaStr := strings.TrimSpace(r.FormValue("fecha"))
	if fechaStr == "" {
		result.Errors["fecha"] = "La fecha es obligatoria"
		result.Success = false
	} else {
		// Parsear la fecha en la zona horaria de Perú (UTC-5)
		loc, err := time.LoadLocation("America/Lima")
		if err != nil {
			// Fallback si no encuentra la zona horaria
			loc = time.FixedZone("UTC-5", -5*60*60)
		}
		
		// Parsear la fecha
		fecha, err := time.ParseInLocation("2006-01-02", fechaStr, loc)
		if err != nil {
			result.Errors["fecha"] = "Formato de fecha inválido"
			result.Success = false
		} else {
			// Ajustar a las 12:00 del mediodía para evitar problemas de zona horaria
			fecha = time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 12, 0, 0, 0, loc)
			
			// Verificar que la FECHA (solo día) no sea futura
			// Comparamos solo las fechas, ignorando las horas
			now := time.Now().In(loc)
			hoy := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			fechaSolo := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, loc)
			
			if fechaSolo.After(hoy) {
				result.Errors["fecha"] = "La fecha no puede ser futura"
				result.Success = false
			} else {
				result.CleanData["fecha"] = fecha
			}
		}
	}

	return result
}

func (v *MovimientoValidator) ValidateUpdate(r *http.Request) ValidationResult {
	result := v.Validate(r)

	// Validar ID del movimiento
	idStr := strings.TrimSpace(r.FormValue("idMovimiento"))
	if idStr == "" {
		result.Errors["idMovimiento"] = "El ID del movimiento es obligatorio"
		result.Success = false
	} else {
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			result.Errors["idMovimiento"] = "ID de movimiento inválido"
			result.Success = false
		} else {
			result.CleanData["idMovimiento"] = id
		}
	}

	return result
}