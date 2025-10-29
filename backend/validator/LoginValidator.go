package validator

import "net/http"

type LoginValidator struct {
	ValidatorBase
}

var LoginValidatorInstance = &LoginValidator{
	ValidatorBase: ValidatorBase{
		Rules: map[string]Rule{
			"email": {
				Required: true,
				Email:    true,
			},
			"password": {
				Required:  true,
				MinLength: 6,
			},
		},
	},
}

func (v *LoginValidator) Validate(r *http.Request) ValidationResult {
	return v.ValidatorBase.Validate(r)
}
