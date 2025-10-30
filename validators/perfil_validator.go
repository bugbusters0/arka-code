package validators

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type PerfilValidator struct{}

var PerfilValidatorInstance = &PerfilValidator{}

// ValidateCrearMiembro valida la creación de un nuevo miembro
func (v *PerfilValidator) ValidateCrearMiembro(r *http.Request) ValidationResult {
	result := NewValidationResult()

	// Validar nombre personal (solo letras y espacios)
	nombrePersonal := strings.TrimSpace(r.FormValue("nombrePersonal"))
	if nombrePersonal == "" {
		result.Errors["nombrePersonal"] = "El nombre es obligatorio"
		result.Success = false
	} else if len(nombrePersonal) < 2 || len(nombrePersonal) > 100 {
		result.Errors["nombrePersonal"] = "El nombre debe tener entre 2 y 100 caracteres"
		result.Success = false
	} else {
		// Solo letras (incluyendo acentos) y espacios
		nombreRegex := `^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`
		matched, _ := regexp.MatchString(nombreRegex, nombrePersonal)
		if !matched {
			result.Errors["nombrePersonal"] = "El nombre solo debe contener letras y espacios"
			result.Success = false
		} else {
			result.CleanData["nombrePersonal"] = nombrePersonal
		}
	}

	// Validar nombre de usuario (alfanumérico y guiones bajos)
	nombreUsuario := strings.TrimSpace(r.FormValue("nombreUsuario"))
	if nombreUsuario == "" {
		result.Errors["nombreUsuario"] = "El nombre de usuario es obligatorio"
		result.Success = false
	} else if len(nombreUsuario) < 3 || len(nombreUsuario) > 30 {
		result.Errors["nombreUsuario"] = "El nombre de usuario debe tener entre 3 y 30 caracteres"
		result.Success = false
	} else {
		usuarioRegex := `^[a-zA-Z0-9_]+$`
		matched, _ := regexp.MatchString(usuarioRegex, nombreUsuario)
		if !matched {
			result.Errors["nombreUsuario"] = "El nombre de usuario solo puede contener letras, números y guiones bajos"
			result.Success = false
		} else {
			result.CleanData["nombreUsuario"] = nombreUsuario
		}
	}

	// Validar rol
	rolStr := strings.TrimSpace(r.FormValue("rol"))
	if rolStr == "" {
		result.Errors["rol"] = "El rol es obligatorio"
		result.Success = false
	} else {
		rol, err := strconv.Atoi(rolStr)
		if err != nil || (rol != 0 && rol != 1) {
			result.Errors["rol"] = "Rol inválido (debe ser 0 o 1)"
			result.Success = false
		} else {
			result.CleanData["rol"] = int8(rol)
		}
	}

	// Validar contraseña (opcional, pero si se proporciona debe cumplir requisitos)
	contrasena := r.FormValue("contrasena")
	if contrasena != "" {
		if len(contrasena) < 5 {
			result.Errors["contrasena"] = "La contraseña debe tener al menos 5 caracteres"
			result.Success = false
		} else {
			// Al menos 1 mayúscula, 1 número, 1 carácter especial
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(contrasena)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(contrasena)
			hasSpecial := regexp.MustCompile(`[.,:\[\]+./#$%&/()=?¿¡!]`).MatchString(contrasena)

			if !hasUpper {
				result.Errors["contrasena"] = "La contraseña debe tener al menos una mayúscula"
				result.Success = false
			} else if !hasNumber {
				result.Errors["contrasena"] = "La contraseña debe tener al menos un número"
				result.Success = false
			} else if !hasSpecial {
				result.Errors["contrasena"] = "La contraseña debe tener al menos un carácter especial (.,:[]+./#$%&/()=?¿¡!)"
				result.Success = false
			} else {
				result.CleanData["contrasena"] = contrasena
			}
		}
	}

	return result
}

// ValidateEditarMiembro valida la edición de un miembro
func (v *PerfilValidator) ValidateEditarMiembro(r *http.Request) ValidationResult {
	result := NewValidationResult()

	// Validar nombre de usuario actual (debe existir)
	nombreUsuarioActual := strings.TrimSpace(r.FormValue("nombreUsuarioActual"))
	if nombreUsuarioActual == "" {
		result.Errors["nombreUsuarioActual"] = "Se requiere el nombre de usuario actual"
		result.Success = false
	} else {
		result.CleanData["nombreUsuarioActual"] = nombreUsuarioActual
	}

	// Validar nuevo nombre personal
	nombrePersonal := strings.TrimSpace(r.FormValue("nombrePersonal"))
	if nombrePersonal != "" {
		if len(nombrePersonal) < 2 || len(nombrePersonal) > 100 {
			result.Errors["nombrePersonal"] = "El nombre debe tener entre 2 y 100 caracteres"
			result.Success = false
		} else {
			nombreRegex := `^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]+$`
			matched, _ := regexp.MatchString(nombreRegex, nombrePersonal)
			if !matched {
				result.Errors["nombrePersonal"] = "El nombre solo debe contener letras y espacios"
				result.Success = false
			} else {
				result.CleanData["nombrePersonal"] = nombrePersonal
			}
		}
	}

	// Validar nuevo nombre de usuario (si se proporciona)
	nombreUsuarioNuevo := strings.TrimSpace(r.FormValue("nombreUsuario"))
	if nombreUsuarioNuevo != "" && nombreUsuarioNuevo != nombreUsuarioActual {
		if len(nombreUsuarioNuevo) < 3 || len(nombreUsuarioNuevo) > 30 {
			result.Errors["nombreUsuario"] = "El nombre de usuario debe tener entre 3 y 30 caracteres"
			result.Success = false
		} else {
			usuarioRegex := `^[a-zA-Z0-9_]+$`
			matched, _ := regexp.MatchString(usuarioRegex, nombreUsuarioNuevo)
			if !matched {
				result.Errors["nombreUsuario"] = "El nombre de usuario solo puede contener letras, números y guiones bajos"
				result.Success = false
			} else {
				result.CleanData["nombreUsuario"] = nombreUsuarioNuevo
			}
		}
	}

	// Validar rol (si se proporciona)
	rolStr := strings.TrimSpace(r.FormValue("rol"))
	if rolStr != "" {
		rol, err := strconv.Atoi(rolStr)
		if err != nil || (rol != 0 && rol != 1) {
			result.Errors["rol"] = "Rol inválido (debe ser 0 o 1)"
			result.Success = false
		} else {
			result.CleanData["rol"] = int8(rol)
		}
	}

	// Validar nueva contraseña (opcional)
	contrasena := r.FormValue("contrasena")
	if contrasena != "" {
		if len(contrasena) < 5 {
			result.Errors["contrasena"] = "La contraseña debe tener al menos 5 caracteres"
			result.Success = false
		} else {
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(contrasena)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(contrasena)
			hasSpecial := regexp.MustCompile(`[.,:\[\]+./#$%&/()=?¿¡!]`).MatchString(contrasena)

			if !hasUpper {
				result.Errors["contrasena"] = "La contraseña debe tener al menos una mayúscula"
				result.Success = false
			} else if !hasNumber {
				result.Errors["contrasena"] = "La contraseña debe tener al menos un número"
				result.Success = false
			} else if !hasSpecial {
				result.Errors["contrasena"] = "La contraseña debe tener al menos un carácter especial"
				result.Success = false
			} else {
				result.CleanData["contrasena"] = contrasena
			}
		}
	}

	return result
}

// ValidateLimite valida la creación/edición de un límite de gasto
func (v *PerfilValidator) ValidateLimite(r *http.Request) ValidationResult {
	result := NewValidationResult()

	// Validar concepto
	nombreConcepto := strings.TrimSpace(r.FormValue("nombreConcepto"))
	if nombreConcepto == "" {
		result.Errors["nombreConcepto"] = "El concepto es obligatorio"
		result.Success = false
	} else {
		result.CleanData["nombreConcepto"] = nombreConcepto
	}

	// Validar periodo
	periodo := strings.TrimSpace(r.FormValue("periodo"))
	if periodo == "" {
		result.Errors["periodo"] = "El periodo es obligatorio"
		result.Success = false
	} else if periodo != "diario" && periodo != "semanal" && periodo != "mensual" && periodo != "anual" {
		result.Errors["periodo"] = "Periodo inválido (debe ser: diario, semanal, mensual o anual)"
		result.Success = false
	} else {
		result.CleanData["periodo"] = periodo
	}

	// Validar límite (debe ser un número positivo)
	limiteStr := strings.TrimSpace(r.FormValue("limite"))
	if limiteStr == "" {
		result.Errors["limite"] = "El límite es obligatorio"
		result.Success = false
	} else {
		limite, err := strconv.ParseFloat(limiteStr, 64)
		if err != nil || limite <= 0 {
			result.Errors["limite"] = "El límite debe ser un número mayor a 0"
			result.Success = false
		} else {
			result.CleanData["limite"] = limite
		}
	}

	// Validar nombre de usuario (a quien se le aplica el límite)
	nombreUsuario := strings.TrimSpace(r.FormValue("nombreUsuario"))
	if nombreUsuario == "" {
		result.Errors["nombreUsuario"] = "El usuario es obligatorio"
		result.Success = false
	} else {
		result.CleanData["nombreUsuario"] = nombreUsuario
	}

	return result
}
