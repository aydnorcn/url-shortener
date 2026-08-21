package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var messages = map[string]string{
	"required":     "The field is required",
	"url":          "Enter a valid URL",
	"shortcode":    "Invalid short code",
	"email":        "Enter a valid email address",
	"numeric":      "Enter a valid number",
	"alpha":        "Enter a valid alpha number",
	"alphanumeric": "Enter a valid alphanumeric number",
	"password":     "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character (. ! or _).",
}

func GetErrorMessage(fe validator.FieldError) string {
	tag := fe.Tag()

	if msg, exists := messages[tag]; exists {
		return msg
	}

	switch tag {
	case "min":
		return fmt.Sprintf("Must be at least %s characters long.", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long.", fe.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters long.", fe.Param())
	case "eqfield":
		return "Values must match."
	default:
		return "Invalid value."
	}
}
