package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	shortCodeRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
	hasUpper       = regexp.MustCompile(`[A-Z]`)
	hasLower       = regexp.MustCompile(`[a-z]`)
	hasNumber      = regexp.MustCompile(`[0-9]`)
	hasSpecial     = regexp.MustCompile(`[.!_]`)
)

func registerCustomValidators(v *validator.Validate) {
	_ = v.RegisterValidation("necessary", validateNecessary)
	_ = v.RegisterValidation("shortcode", validateShortCode)
	_ = v.RegisterValidation("password", validatePassword)
}

// validateNecessary checks if a field is provided and not empty/whitespace.
func validateNecessary(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return strings.TrimSpace(field.String()) != ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface:
		if field.IsNil() {
			return false
		}
		elem := field.Elem()
		if elem.Kind() == reflect.String {
			return strings.TrimSpace(elem.String()) != ""
		}
		return !elem.IsZero()
	default:
		return field.IsValid() && !field.IsZero()
	}
}

func validateShortCode(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "" {
		return true
	}
	return shortCodeRegex.MatchString(val)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true
	}
	return len(password) >= 8 &&
		hasUpper.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasNumber.MatchString(password) &&
		hasSpecial.MatchString(password)
}
