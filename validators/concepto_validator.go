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
	if periodoTipo != "" && periodoTipo != "diario" && periodoTipo != "quincenal" && periodoTipo != "mensual" {
		result.Errors["periodo_tipo"] = "El tipo de período debe ser: diario, quincenal o mensual"
		result.Success = false
	} else {
		result.CleanData["periodo_tipo"] = periodoTipo
	}
	// Dia Periodo tipo para desembolso
	diaDesembolsoStr := strings.TrimSpace(r.FormValue("dia_desembolso_planejado"))
	if diaDesembolsoStr != "" {
		diaDesembolso, err := strconv.ParseInt(diaDesembolsoStr, 10, 8)
		if err != nil {
			result.Errors["dia_desembolso_planejado"] = "El día debe ser un número válido (Desembolso planificado)."
			result.Success = false
		} else if diaDesembolso < 1 || diaDesembolso > 31 {
			result.Errors["dia_desembolso_planejado"] = "El día debe estar en un rango de 1 - 31"
			result.Success = false
		} else {
			// GUARDAR EL VALOR CONVERTIDO como int8
			result.CleanData["dia_desembolso_planejado"] = int8(diaDesembolso)
		}
	} else {
		// Si no hay valor, establecer nil para que no se guarde en BD
		result.CleanData["dia_desembolso_planejado"] = nil
	}

	// Validación adicional: si período es mensual, requiere día
	if periodoTipo == "mensual" && diaDesembolsoStr == "" {
		result.Errors["dia_desembolso_planejado"] = "El día es requerido para período mensual"
		result.Success = false
	}

	// Si período es quincenal, forzar día 15
	if periodoTipo == "quincenal" {
		result.CleanData["dia_desembolso_planejado"] = int8(15)
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
	if limiteTipo != "" && limiteTipo != "diario" && limiteTipo != "quincenal" && limiteTipo != "mensual" {
		result.Errors["limite_tipo"] = "El tipo de límite debe ser: diario, quincenal o mensual"
		result.Success = false
	} else {
		result.CleanData["limite_tipo"] = limiteTipo
	}

	// Dia Periodo tipo para desembolso
	diaLimiteStr := strings.TrimSpace(r.FormValue("dia_limite_tipo"))
	if diaLimiteStr != "" {
		diaLimite, err := strconv.ParseInt(diaLimiteStr, 10, 8)
		if err != nil {
			result.Errors["dia_limite_tipo"] = "El día debe ser un número válido (Establecer límite)."
			result.Success = false
		} else if diaLimite < 1 || diaLimite > 31 {
			result.Errors["dia_limite_tipo"] = "El día debe estar en un rango de 1 - 31"
			result.Success = false
		} else {
			// ✅ GUARDAR EL VALOR CONVERTIDO como int8
			result.CleanData["dia_limite_tipo"] = int8(diaLimite)
		}
	} else {
		// Si no hay valor, establecer nil para que no se guarde en BD
		result.CleanData["dia_limite_tipo"] = nil
	}

	// Validación adicional: si límite tipo es mensual, requiere día
	if limiteTipo == "mensual" && diaLimiteStr == "" {
		result.Errors["dia_limite_tipo"] = "El día es requerido para límite mensual"
		result.Success = false
	}

	// Si límite tipo es quincenal, forzar día 15
	if limiteTipo == "quincenal" {
		result.CleanData["dia_limite_tipo"] = int8(15)
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
