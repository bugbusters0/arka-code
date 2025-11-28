package validators

import (
	"fmt"
	"net/http"
	"strings"
)

type RegisterValidator struct {
	*ValidatorBase
}

var RegisterValidatorInstance = &RegisterValidator{}

func (v *RegisterValidator) Validate(r *http.Request) ValidationResult {
	// Validación base (administrador)
	validator := &ValidatorBase{
		Email:           strings.TrimSpace(r.FormValue("email")),
		Telefono:        strings.TrimSpace(r.FormValue("telefono")),
		ConfirmPassword: r.FormValue("confirmPassword"),
		Nombre:          strings.TrimSpace(r.FormValue("adminNombre")),
		Contrasena:      r.FormValue("password"),
	}

	fieldsToValidate := map[string]string{
		"email":           validator.Email,
		"nombre":          validator.Nombre,
		"contrasena":      validator.Contrasena,
		"confirmPassword": validator.ConfirmPassword,
		"telefono":        validator.Telefono,
		"contraPersonal":  r.FormValue("contraPersonal"),
		"nombreUsuario":   strings.TrimSpace(r.FormValue("nombreUsuario")),
	}

	result := validator.ValidateFields(fieldsToValidate)

	// ✅ Validar contraseñas
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	if password != confirmPassword {
		result.Errors["confirmPassword"] = "Las contraseñas no coinciden"
		result.Success = false
	}

	// ✅ Validar contraPersonal (PIN)
	contraPersonal := r.FormValue("contraPersonal")
	if len(contraPersonal) < 4 || len(contraPersonal) > 6 {
		result.Errors["contraPersonal"] = "La contraseña personal debe tener entre 4 y 6 dígitos"
		result.Success = false
	}

	// ✅ Validar nombreUsuario (único)
	nombreUsuario := strings.TrimSpace(r.FormValue("nombreUsuario"))
	if nombreUsuario == "" {
		result.Errors["nombreUsuario"] = "El nombre de usuario es requerido"
		result.Success = false
	} else if len(nombreUsuario) < 3 {
		result.Errors["nombreUsuario"] = "El nombre de usuario debe tener al menos 3 caracteres"
		result.Success = false
	}

	// ✅ NUEVO: Validar miembros adicionales
	memberErrors := v.validateMembers(r)
	for field, errorMsg := range memberErrors {
		result.Errors[field] = errorMsg
	}
	if len(memberErrors) > 0 {
		result.Success = false
	}

	return result
}

// ✅ NUEVA FUNCIÓN: Validar miembros adicionales
func (v *RegisterValidator) validateMembers(r *http.Request) map[string]string {
	errors := make(map[string]string)

	miembroNombres := r.Form["miembroNombre[]"]
	miembroUsuarios := r.Form["miembroUsuario[]"]
	miembroContras := r.Form["miembroContra[]"]

	// Si no hay miembros, no hay errores
	if len(miembroNombres) == 0 {
		return errors
	}

	// Validar cada miembro
	baseValidator := &ValidatorBase{}
	usuarioNames := make(map[string]bool) // Para verificar duplicados

	for i := 0; i < len(miembroNombres); i++ {
		// Si todos los campos están vacíos, es un miembro vacío (skip)
		if miembroNombres[i] == "" && miembroUsuarios[i] == "" && miembroContras[i] == "" {
			continue
		}

		// Validar que todos los campos estén completos
		if miembroNombres[i] == "" {
			errors[fmt.Sprintf("miembroNombre_%d", i)] = "El nombre del miembro es requerido"
		} else if !baseValidator.IsValidNombre(miembroNombres[i]) {
			errors[fmt.Sprintf("miembroNombre_%d", i)] = "Nombre del miembro inválido (solo letras y espacios, 2-100 caracteres)"
		}

		if miembroUsuarios[i] == "" {
			errors[fmt.Sprintf("miembroUsuario_%d", i)] = "El nombre de usuario del miembro es requerido"
		} else if len(miembroUsuarios[i]) < 3 {
			errors[fmt.Sprintf("miembroUsuario_%d", i)] = "El nombre de usuario debe tener al menos 3 caracteres"
		} else if usuarioNames[miembroUsuarios[i]] {
			errors[fmt.Sprintf("miembroUsuario_%d", i)] = "Nombre de usuario duplicado en los miembros"
		} else {
			usuarioNames[miembroUsuarios[i]] = true
		}

		if miembroContras[i] == "" {
			errors[fmt.Sprintf("miembroContra_%d", i)] = "La contraseña personal del miembro es requerida"
		} else if len(miembroContras[i]) < 4 || len(miembroContras[i]) > 6 {
			errors[fmt.Sprintf("miembroContra_%d", i)] = "La contraseña personal debe tener entre 4 y 6 dígitos"
		}
	}

	return errors
}
