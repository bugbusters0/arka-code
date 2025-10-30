package validators

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ConceptoValidator struct {
	*ValidatorBase
}

var ConceptoValidatorInstance = &ConceptoValidator{}

func (v *ConceptoValidator) Validate(r *http.Request) ValidationResult {
	validator := &ValidatorBase{
		Nombre: strings.TrimSpace(r.FormValue("nombre")),
		Color:  strings.TrimSpace(r.FormValue("color")),
	}

	fieldsToValidate := map[string]string{
		"nombre": validator.Nombre,
		"color":  validator.Color,
	}

	result := validator.ValidateFields(fieldsToValidate)

	// ✅ Validaciones específicas para concepto
	v.validateConceptoFields(r, &result)

	return result
}

// validateConceptoFields valida campos específicos de concepto
func (v *ConceptoValidator) validateConceptoFields(r *http.Request, result *ValidationResult) {
	// Validar tipo (gasto/ingreso)
	tipo := strings.TrimSpace(r.FormValue("tipo"))
	if tipo == "" {
		result.Errors["tipo"] = "El tipo de concepto es requerido"
		result.Success = false
	} else if tipo != "gasto" && tipo != "ingreso" {
		result.Errors["tipo"] = "El tipo debe ser 'gasto' o 'ingreso'"
		result.Success = false
	} else {
		result.CleanData["tipo"] = tipo
	}

	// Validar id_icono
	idIcono := strings.TrimSpace(r.FormValue("id_icono"))
	if idIcono == "" {
		result.Errors["id_icono"] = "Debe seleccionar un ícono"
		result.Success = false
	} else {
		iconoID, err := strconv.Atoi(idIcono)
		if err != nil || iconoID <= 0 {
			result.Errors["id_icono"] = "El ícono seleccionado no es válido"
			result.Success = false
		} else {
			result.CleanData["id_icono"] = iconoID
		}
	}

	// Validar nombre (longitud específica para BD)
	nombre := strings.TrimSpace(r.FormValue("nombre"))
	if len(nombre) > 40 {
		result.Errors["nombre"] = "El nombre no puede tener más de 40 caracteres"
		result.Success = false
	}

	// Validar configuracion_usuario
	configUsuario := strings.TrimSpace(r.FormValue("configuracion_usuario"))
	if configUsuario != "true" && configUsuario != "false" {
		result.Errors["configuracion_usuario"] = "Configuración de usuario inválida"
		result.Success = false
	} else {
		result.CleanData["configuracion_usuario"] = configUsuario == "true"
	}

	// ✅ Validar campos de personalización (opcionales)
	v.validatePersonalizacion(r, result)

	// Si hay errores, actualizar el estado
	if len(result.Errors) > 0 {
		result.Success = false
	}
}

// validatePersonalizacion valida los campos opcionales de personalización
func (v *ConceptoValidator) validatePersonalizacion(r *http.Request, result *ValidationResult) {
	// Desembolso planificado
	desembolsoStr := strings.TrimSpace(r.FormValue("desembolso_planejado"))
	if desembolsoStr != "" {
		desembolso, err := strconv.ParseFloat(desembolsoStr, 64)
		if err != nil || desembolso < 0 {
			result.Errors["desembolso_planejado"] = "El desembolso planificado debe ser un número positivo"
			result.Success = false
		} else {
			result.CleanData["desembolso_planejado"] = desembolso
		}
	} else {
		result.CleanData["desembolso_planejado"] = 0.0
	}

	// Período tipo para desembolso
	periodoTipo := strings.TrimSpace(r.FormValue("periodo_tipo"))
	if periodoTipo != "" && periodoTipo != "diario" && periodoTipo != "semanal" && periodoTipo != "mensual" {
		result.Errors["periodo_tipo"] = "El tipo de período debe ser: diario, semanal o mensual"
		result.Success = false
	} else {
		result.CleanData["periodo_tipo"] = periodoTipo
	}

	// Límite de monto
	limiteStr := strings.TrimSpace(r.FormValue("limite_monto"))
	if limiteStr != "" {
		limite, err := strconv.ParseFloat(limiteStr, 64)
		if err != nil || limite < 0 {
			result.Errors["limite_monto"] = "El límite debe ser un número positivo"
			result.Success = false
		} else {
			result.CleanData["limite_monto"] = limite
		}
	} else {
		result.CleanData["limite_monto"] = 0.0
	}

	// Límite tipo
	limiteTipo := strings.TrimSpace(r.FormValue("limite_tipo"))
	if limiteTipo != "" && limiteTipo != "diario" && limiteTipo != "semanal" && limiteTipo != "mensual" {
		result.Errors["limite_tipo"] = "El tipo de límite debe ser: diario, semanal o mensual"
		result.Success = false
	} else {
		result.CleanData["limite_tipo"] = limiteTipo
	}

	// Día del período (opcional, 1-31)
	diaPeriodoStr := strings.TrimSpace(r.FormValue("dia_periodo_planificado"))
	if diaPeriodoStr != "" {
		diaPeriodo, err := strconv.Atoi(diaPeriodoStr)
		if err != nil || diaPeriodo < 1 || diaPeriodo > 31 {
			result.Errors["dia_periodo_planificado"] = "El día del período debe estar entre 1 y 31"
			result.Success = false
		} else {
			result.CleanData["dia_periodo_planificado"] = diaPeriodo
		}
	}

	// Notificación (booleano)
	notificacionStr := strings.TrimSpace(r.FormValue("notificacion"))
	if notificacionStr != "" {
		notificacion, err := strconv.ParseBool(notificacionStr)
		if err != nil {
			result.Errors["notificacion"] = "La notificación debe ser verdadero o falso"
			result.Success = false
		} else {
			result.CleanData["notificacion"] = notificacion
		}
	} else {
		result.CleanData["notificacion"] = false
	}
}

// ValidateUpdate valida la actualización de un concepto
func (v *ConceptoValidator) ValidateUpdate(r *http.Request) ValidationResult {
	result := v.Validate(r)

	// Validaciones adicionales para actualización
	nombreConcepto := strings.TrimSpace(r.FormValue("nombre_original"))
	if nombreConcepto == "" {
		result.Errors["nombre_original"] = "El nombre original del concepto es requerido para actualizar"
		result.Success = false
	} else {
		result.CleanData["nombre_original"] = nombreConcepto
	}

	return result
}

// ValidatePersonalizacion valida solo los campos de personalización
func (v *ConceptoValidator) ValidatePersonalizacion(r *http.Request) ValidationResult {
	result := NewValidationResult()

	// Validar que el concepto existe
	nombreConcepto := strings.TrimSpace(r.FormValue("nombre_concepto"))
	correoFamilia := strings.TrimSpace(r.FormValue("correo_familia"))

	if nombreConcepto == "" {
		result.Errors["nombre_concepto"] = "El nombre del concepto es requerido"
	}
	if correoFamilia == "" {
		result.Errors["correo_familia"] = "El correo de la familia es requerido"
	}

	if len(result.Errors) > 0 {
		result.Success = false
		return result
	}

	result.CleanData["nombre_concepto"] = nombreConcepto
	result.CleanData["correo_familia"] = correoFamilia

	// Validar campos de personalización
	v.validatePersonalizacion(r, &result)

	return result
}

// ValidateNombreUnico verifica si el nombre del concepto es único para la familia
// Esta función se usaría en el controlador después de validar los datos básicos
func (v *ConceptoValidator) ValidateNombreUnico(nombreConcepto, correoFamilia string, esActualizacion bool, nombreOriginal string) (bool, string) {
	// Esta validación requiere acceso a la base de datos
	// Se implementaría en el controlador llamando al modelo

	if esActualizacion && nombreConcepto == nombreOriginal {
		return true, "" // Mismo nombre en actualización, no hay problema
	}

	// Para implementación completa, necesitarías inyectar el modelo o hacer esta validación en el controlador
	return true, "" // Placeholder
}

// Helper para validar formato de ícono (FontAwesome)
func (v *ConceptoValidator) IsValidIcono(icono string) bool {
	if icono == "" {
		return false
	}

	// Validar formato básico de clase FontAwesome
	// Ej: "fa-solid fa-house", "fa-regular fa-user", etc.
	iconoRegex := `^fa-[a-z]+(-[a-z]+)* fa-[a-z]+(-[a-z]+)*$`
	matched, _ := regexp.MatchString(iconoRegex, strings.TrimSpace(icono))
	return matched
}
