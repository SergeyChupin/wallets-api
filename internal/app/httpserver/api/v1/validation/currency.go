package validation

import (
	"github.com/go-playground/validator/v10"
)

var (
	currencies = map[string]struct{}{
		"USD": {},
	}

	CurrencyValidator = func(fl validator.FieldLevel) bool {
		return hasCurrency(fl.Field().String())
	}
)

func hasCurrency(currency string) bool {
	_, ok := currencies[currency]
	return ok
}
