package validators

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidatorBase - Estructura base para validadores
type ValidatorBase struct {
	Email           string
	Nombre          string
	Numero          int
	Contrasena      string
	Telefono        string
	Color           string
	Rol             string
	FechaNac        string
	ConfirmPassword string
}

// NewValidatorBase crea una nueva instancia del validador base
func NewValidatorBase(r *http.Request) *ValidatorBase {
	return &ValidatorBase{
		Email:           strings.TrimSpace(r.FormValue("email")),
		Nombre:          strings.TrimSpace(r.FormValue("nombre")),
		Numero:          stringToInt(r.FormValue("numero")),
		Contrasena:      r.FormValue("contrasena"),
		Telefono:        strings.TrimSpace(r.FormValue("telefono")),
		Color:           strings.TrimSpace(r.FormValue("color")),
		Rol:             strings.TrimSpace(r.FormValue("rol")),
		FechaNac:        strings.TrimSpace(r.FormValue("fecha_nac")),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
}

// ValidationResult resultado de validación
type ValidationResult struct {
	Success   bool                   `json:"success"`
	Errors    map[string]string      `json:"errors"`
	CleanData map[string]interface{} `json:"cleanData"`
}

// NewValidationResult crea un nuevo resultado de validación
func NewValidationResult() ValidationResult {
	return ValidationResult{
		Success:   true,
		Errors:    make(map[string]string),
		CleanData: make(map[string]interface{}),
	}
}

// ========== MÉTODOS DE VALIDACIÓN ==========

// IsValidEmail valida formato de email
func (v *ValidatorBase) IsValidEmail(email string) bool {
	emailToCheck := email
	if emailToCheck == "" {
		emailToCheck = v.Email
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, emailToCheck)
	return matched
}

// IsValidNombre valida nombre (2-100 chars, solo letras y espacios)
func (v *ValidatorBase) IsValidNombre(nombre string) bool {
	nombreToCheck := nombre
	if nombreToCheck == "" {
		nombreToCheck = v.Nombre
	}

	if len(strings.TrimSpace(nombreToCheck)) < 2 || len(strings.TrimSpace(nombreToCheck)) > 100 {
		return false
	}

	// Solo letras y espacios
	nombreRegex := `^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`
	matched, _ := regexp.MatchString(nombreRegex, strings.TrimSpace(nombreToCheck))
	return matched
}

// IsValidNumero valida número positivo
func (v *ValidatorBase) IsValidNumero(numero int) bool {
	numeroToCheck := numero
	if numeroToCheck == 0 {
		numeroToCheck = v.Numero
	}
	return numeroToCheck > 0
}

// IsValidContrasena valida contraseña (mínimo 6 chars, mayúscula, minúscula, número)
func (v *ValidatorBase) IsValidContrasena(contrasena string) bool {
	contrasenaToCheck := contrasena
	if contrasenaToCheck == "" {
		contrasenaToCheck = v.Contrasena
	}

	if len(contrasenaToCheck) < 6 {
		return false
	}

	// Al menos una mayúscula, una minúscula y un número
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(contrasenaToCheck)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(contrasenaToCheck)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(contrasenaToCheck)

	return hasUpper && hasLower && hasNumber
}

// IsValidTelefono valida teléfono (9 dígitos para Perú)
func (v *ValidatorBase) IsValidTelefono(telefono string) bool {
	telefonoToCheck := telefono
	if telefonoToCheck == "" {
		telefonoToCheck = v.Telefono
	}

	// 9 dígitos sin espacios
	telefonoRegex := `^[0-9]{9}$`
	matched, _ := regexp.MatchString(telefonoRegex, telefonoToCheck)
	return matched
}

// IsValidColor valida formato hexadecimal de color
func (v *ValidatorBase) IsValidColor(color string) bool {
	colorToCheck := color
	if colorToCheck == "" {
		colorToCheck = v.Color
	}

	colorRegex := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	matched, _ := regexp.MatchString(colorRegex, colorToCheck)
	return matched
}

// IsValidRol valida rol (admin o miembro)
func (v *ValidatorBase) IsValidRol(rol string) bool {
	rolToCheck := rol
	if rolToCheck == "" {
		rolToCheck = v.Rol
	}

	rolesValidos := map[string]bool{
		"admin":   true,
		"miembro": true,
	}

	return rolesValidos[strings.ToLower(strings.TrimSpace(rolToCheck))]
}

// IsValidFechaNacimiento valida fecha de nacimiento
func (v *ValidatorBase) IsValidFechaNacimiento(fecha string) bool {
	fechaToCheck := fecha
	if fechaToCheck == "" {
		fechaToCheck = v.FechaNac
	}

	// Validar formato YYYY-MM-DD
	_, err := time.Parse("2006-01-02", fechaToCheck)
	if err != nil {
		return false
	}

	// No puede ser futura
	fechaObj, _ := time.Parse("2006-01-02", fechaToCheck)
	if fechaObj.After(time.Now()) {
		return false
	}

	// Debe tener al menos 5 años
	edad := time.Now().Year() - fechaObj.Year()
	return edad >= 5
}

// IsValidContraPersonal valida contraseña personal (PIN 4-6 chars)
func (v *ValidatorBase) IsValidContraPersonal(contra string) bool {
	contraToCheck := contra
	if contraToCheck == "" {
		contraToCheck = v.Contrasena
	}

	length := len(contraToCheck)
	return length >= 4 && length <= 6
}

// PasswordsMatch verifica que contraseñas coincidan
func (v *ValidatorBase) PasswordsMatch(password, confirmPassword string) bool {
	if password == "" {
		password = v.Contrasena
	}
	if confirmPassword == "" {
		confirmPassword = v.ConfirmPassword
	}

	return password == confirmPassword
}

// ========== MÉTODOS HELPER ==========

// ValidateFields valida múltiples campos y devuelve errores
func (v *ValidatorBase) ValidateFields(fields map[string]string) ValidationResult {
	result := NewValidationResult()

	for field, value := range fields {
		// Determinar qué método de validación usar
		isValid, cleanValue := v.validateField(field, value)

		if !isValid {
			result.Errors[field] = v.getErrorMessage(field)
		} else {
			result.CleanData[field] = cleanValue
		}
	}

	result.Success = len(result.Errors) == 0
	return result
}

// validateField valida un campo individual
func (v *ValidatorBase) validateField(field, value string) (bool, interface{}) {
	switch field {
	case "email":
		return v.IsValidEmail(value), v.sanitizeEmail(value)
	case "nombre":
		return v.IsValidNombre(value), v.sanitizeString(value)
	case "numero":
		num := stringToInt(value)
		return v.IsValidNumero(num), num
	case "contrasena":
		return v.IsValidContrasena(value), value // No sanitizar contraseñas
	case "telefono":
		return v.IsValidTelefono(value), v.sanitizeInt(value)
	case "color":
		return v.IsValidColor(value), v.sanitizeString(value)
	case "rol":
		return v.IsValidRol(value), v.sanitizeString(value)
	case "fecha_nac":
		return v.IsValidFechaNacimiento(value), v.sanitizeString(value)
	case "contraPersonal":
		return v.IsValidContraPersonal(value), value // No sanitizar
	case "confirmPassword":
		return true, value // No sanitizar, validar con PasswordsMatch después
	default:
		// Si no hay validador específico, solo sanitizar
		return true, v.sanitizeString(value)
	}
}

// getErrorMessage retorna mensajes de error personalizados
func (v *ValidatorBase) getErrorMessage(field string) string {
	messages := map[string]string{
		"email":          "El correo electrónico no es válido",
		"nombre":         "El nombre debe tener entre 2 y 100 caracteres y solo letras",
		"contrasena":     "La contraseña debe tener al menos 6 caracteres, una mayúscula, una minúscula y un número",
		"telefono":       "El teléfono debe tener 9 dígitos",
		"color":          "El color debe estar en formato hexadecimal (#RRGGBB)",
		"rol":            "El rol debe ser 'admin' o 'miembro'",
		"fecha_nac":      "La fecha de nacimiento no es válida",
		"contraPersonal": "La contraseña personal debe tener entre 4 y 6 caracteres",
		"numero":         "El número debe ser mayor a 0",
	}

	if msg, exists := messages[field]; exists {
		return msg
	}
	return "Campo " + field + " inválido"
}

// ========== MÉTODOS DE SANITIZACIÓN ==========

func (v *ValidatorBase) sanitizeString(value string) string {
	return strings.TrimSpace(value)
}

func (v *ValidatorBase) sanitizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (v *ValidatorBase) sanitizeInt(value string) int {
	return stringToInt(value)
}

// ========== FUNCIONES HELPER ==========

func stringToInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
