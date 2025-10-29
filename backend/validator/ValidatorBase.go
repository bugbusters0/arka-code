package validator

import "net/http"

type Validator interface {
	Validate(r *http.Request) ValidationResult
}

type ValidatorBase struct {
	Rules map[string]Rule
}

type Rule struct {
	Required  bool
	MinLength int
	MaxLength int
	Email     bool
	Custom    func(value string) error
}

func (vb *ValidatorBase) Validate(r *http.Request) ValidationResult {
	errors := make(map[string]string)

	for field, rule := range vb.Rules {
		value := r.FormValue(field)

		if rule.Required && value == "" {
			errors[field] = "Este campo es requerido"
			continue
		}

		if rule.MinLength > 0 && len(value) < rule.MinLength {
			errors[field] = "Muy corto"
		}

		if rule.Email && !isValidEmail(value) {
			errors[field] = "Email inválido"
		}

		if rule.Custom != nil {
			if err := rule.Custom(value); err != nil {
				errors[field] = err.Error()
			}
		}
	}

	return ValidationResult{
		Success: len(errors) == 0,
		Errors:  errors,
	}
}
