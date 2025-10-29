package validator

import (
	"net/http"
	"regexp"
)

type RegisterValidator struct{}

var RegisterValidatorInstance = &RegisterValidator{}

func (v *RegisterValidator) Validate(r *http.Request) ValidationResult {
	errors := make(map[string]string)

	nombre := r.FormValue("nombre")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	familyName := r.FormValue("family_name")

	// Validar nombre
	if nombre == "" {
		errors["nombre"] = "El nombre es requerido"
	} else if len(nombre) < 2 {
		errors["nombre"] = "El nombre debe tener al menos 2 caracteres"
	}

	// Validar email
	if email == "" {
		errors["email"] = "El email es requerido"
	} else if !isValidEmail(email) {
		errors["email"] = "El formato del email no es válido"
	}

	// Validar contraseña
	if password == "" {
		errors["password"] = "La contraseña es requerida"
	} else if len(password) < 6 {
		errors["password"] = "La contraseña debe tener al menos 6 caracteres"
	}

	// Validar confirmación
	if password != confirmPassword {
		errors["confirm_password"] = "Las contraseñas no coinciden"
	}

	// Validar nombre de familia
	if familyName == "" {
		errors["family_name"] = "El nombre de la familia es requerido"
	}

	return ValidationResult{
		Success: len(errors) == 0,
		Errors:  errors,
	}
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}
