package validator

import "net/http"

type ConceptoValidator struct{}

var ConceptoValidatorInstance = &ConceptoValidator{}

func (v *ConceptoValidator) Validate(r *http.Request) ValidationResult {
	errors := make(map[string]string)

	nombre := r.FormValue("nombre")
	tipo := r.FormValue("tipo")

	if nombre == "" {
		errors["nombre"] = "El nombre del concepto es requerido"
	} else if len(nombre) < 2 {
		errors["nombre"] = "El nombre debe tener al menos 2 caracteres"
	}

	if tipo == "" {
		errors["tipo"] = "El tipo es requerido"
	} else if tipo != "gasto" && tipo != "ingreso" {
		errors["tipo"] = "El tipo debe ser 'gasto' o 'ingreso'"
	}

	return ValidationResult{
		Success: len(errors) == 0,
		Errors:  errors,
	}
}
