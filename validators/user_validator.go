package validators

import (
	"net/http"
	"regexp"
)

type UserValidator struct{}

var UserValidatorInstance = &UserValidator{}

func (v *UserValidator) ValidateRegister(r *http.Request) ValidationResult {
	result := NewValidationResult()

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	name := r.FormValue("name")

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		result.Errors["email"] = "El email es requerido"
	} else if !emailRegex.MatchString(email) {
		result.Errors["email"] = "El email no es válido"
	}

	if name == "" {
		result.Errors["name"] = "El nombre es requerido"
	} else if len(name) < 2 {
		result.Errors["name"] = "El nombre debe tener al menos 2 caracteres"
	}

	if password == "" {
		result.Errors["password"] = "La contraseña es requerida"
	} else if len(password) < 6 {
		result.Errors["password"] = "La contraseña debe tener al menos 6 caracteres"
	}

	if password != confirmPassword {
		result.Errors["confirm_password"] = "Las contraseñas no coinciden"
	}

	result.Success = len(result.Errors) == 0
	return result
}

func (v *UserValidator) ValidateLogin(r *http.Request) ValidationResult {
	result := NewValidationResult()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		result.Errors["email"] = "El email es requerido"
	}

	if password == "" {
		result.Errors["password"] = "La contraseña es requerida"
	}

	result.Success = len(result.Errors) == 0
	return result
}
