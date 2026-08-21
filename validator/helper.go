package validator

import (
	"errors"
	"url-shortener/appErrors"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(value any) *appErrors.AppError {
	if value == nil {
		return nil
	}

	v := GetValidator()
	err := v.Struct(value)
	if err == nil {
		return nil
	}

	var valErrors validator.ValidationErrors
	if !errors.As(err, &valErrors) {
		return appErrors.NewBadRequestError("Invalid struct")
	}

	details := make(map[string]string)

	for _, fe := range valErrors {
		fieldName := fe.Field()
		msg := GetErrorMessage(fe)

		if _, exists := details[fieldName]; !exists {
			details[fieldName] = msg
		} else {
			if fe.Tag() == "required" {
				details[fieldName] = msg
			}
		}
	}

	return appErrors.NewValidationError(details)
}
