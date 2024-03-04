package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/stevenysy/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// Check if currency is supported
		return util.IsSupportedCurrency(currency)
	}
	// Field is not string
	return false
}
