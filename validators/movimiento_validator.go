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

/****************************************/
/*           FNVal_Movimiento           */
/****************************************/

// Validate procesa y valida los datos de un formulario HTTP para la creación de un movimiento.
// Propósito: Garantizar que los campos obligatorios (`monto`, `nombreConcepto`, `fecha`)
//
//	cumplan con los requisitos de negocio y formato, y que el campo opcional (`descripcion`)
//	no exceda su límite. Realiza la conversión de tipos y la validación de la lógica de negocio (fecha no futura).
//
// Parámetros:
// - r: *http.Request que contiene los datos del formulario (r.FormValue).
// Retorno: Un objeto `ValidationResult` que contiene el estado de la validación (`Success`),
//
//	los mensajes de error (`Errors`) y los datos limpios y convertidos (`CleanData`).
//
// Uso: Se invoca en el controlador `Crear` antes de intentar crear un movimiento.
func (v *MovimientoValidator) Validate(r *http.Request) ValidationResult {
	// Inicializa el objeto de resultado de validación.
	result := NewValidationResult()

	// ----------------------------------------------------------------------
	// Validar monto
	// ----------------------------------------------------------------------
	montoStr := strings.TrimSpace(r.FormValue("monto"))
	if montoStr == "" {
		// Error si el campo está vacío.
		result.Errors["monto"] = "El monto es obligatorio"
		result.Success = false
	} else {
		// Intenta convertir el string a float64.
		monto, err := strconv.ParseFloat(montoStr, 64)
		if err != nil || monto <= 0 {
			// Error si no es un número válido o si es menor o igual a cero.
			result.Errors["monto"] = "El monto debe ser un número mayor a 0"
			result.Success = false
		} else {
			// Almacena el valor float limpio.
			result.CleanData["monto"] = monto
		}
	}

	// ----------------------------------------------------------------------
	// Validar concepto
	// ----------------------------------------------------------------------
	concepto := strings.TrimSpace(r.FormValue("nombreConcepto"))
	if concepto == "" {
		// Error si el campo está vacío (no se seleccionó un concepto).
		result.Errors["nombreConcepto"] = "Debe seleccionar un concepto"
		result.Success = false
	} else {
		// Almacena el valor string limpio.
		result.CleanData["nombreConcepto"] = concepto
	}

	// ----------------------------------------------------------------------
	// Validar descripción (opcional)
	// ----------------------------------------------------------------------
	descripcion := strings.TrimSpace(r.FormValue("descripcion"))
	if len(descripcion) > 500 {
		// Error si la descripción excede el límite de 500 caracteres.
		result.Errors["descripcion"] = "La descripción no puede exceder 500 caracteres"
		result.Success = false
	} else if descripcion != "" {
		// Si no está vacía y es válida, almacena la descripción.
		result.CleanData["descripcion"] = descripcion
	}

	// ----------------------------------------------------------------------
	// Validar fecha
	// ----------------------------------------------------------------------
	fechaStr := strings.TrimSpace(r.FormValue("fecha"))
	if fechaStr == "" {
		// Error si el campo de fecha está vacío.
		result.Errors["fecha"] = "La fecha es obligatoria"
		result.Success = false
	} else {
		// Intenta cargar la zona horaria de Perú (America/Lima).
		loc, err := time.LoadLocation("America/Lima")
		if err != nil {
			// Fallback si la zona horaria no se encuentra en el sistema.
			loc = time.FixedZone("UTC-5", -5*60*60)
		}

		// Parsear la fecha usando el formato esperado (YYYY-MM-DD) y la zona horaria local.
		fecha, err := time.ParseInLocation("2006-01-02", fechaStr, loc)
		if err != nil {
			// Error si el formato de la fecha es incorrecto.
			result.Errors["fecha"] = "Formato de fecha inválido"
			result.Success = false
		} else {
			// Ajustar la hora a las 12:00 PM en la zona horaria de Lima para asegurar consistencia y evitar
			// problemas de cambio de día al manipular zonas horarias.
			fecha = time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 12, 0, 0, 0, loc)

			// Verificar que la FECHA no sea futura (Lógica de Negocio)
			// Obtiene la hora actual en la zona de Lima.
			now := time.Now().In(loc)
			// Crea un objeto time que representa el inicio del día de hoy (00:00:00).
			hoy := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
			// Crea un objeto time que representa el inicio del día de la fecha validada.
			fechaSolo := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, loc)

			if fechaSolo.After(hoy) {
				// Error si la fecha seleccionada es posterior al día de hoy.
				result.Errors["fecha"] = "La fecha no puede ser futura"
				result.Success = false
			} else {
				// Almacena el valor time.Time limpio y ajustado.
				result.CleanData["fecha"] = fecha
			}
		}
	}

	// Retorna el resultado final de la validación.
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
