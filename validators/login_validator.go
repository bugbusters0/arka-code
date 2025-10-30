package validators

import "net/http"

type LoginValidator struct {
	*ValidatorBase
}

var LoginValidatorInstance = &LoginValidator{}

func (v *LoginValidator) Validate(r *http.Request) ValidationResult {
	// Usar composición con ValidatorBase (equivalente a herencia en PHP)
	validator := &LoginValidator{
		ValidatorBase: NewValidatorBase(r),
	}

	// Campos a validar (igual que en tu PHP)
	fieldsToValidate := map[string]string{
		"email":      r.FormValue("email"),
		"contrasena": r.FormValue("contrasena"),
	}

	// Usar el método de ValidatorBase (equivalente a validateFields en PHP)
	result := validator.ValidatorBase.ValidateFields(fieldsToValidate)

	// Validación adicional específica para login
	if !result.Success {
		return result
	}

	// Puedes agregar más validaciones específicas del login aquí
	return result
}
