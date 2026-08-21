package appErrors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(details map[string]string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed.",
		Details: details,
		Status:  http.StatusBadRequest,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

var (
	ErrURLNotFound = &AppError{
		Code:    "URL_NOT_FOUND",
		Message: "URL not found",
		Status:  http.StatusNotFound,
	}

	ErrUserNotFound = &AppError{
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
		Status:  http.StatusNotFound,
	}

	ErrInvalidURL = &AppError{
		Code:    "INVALID_URL",
		Message: "Invalid URL",
		Status:  http.StatusBadRequest,
	}

	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized",
		Status:  http.StatusUnauthorized,
	}

	ErrServerError = &AppError{
		Code:    "SERVER_ERROR",
		Message: "Server error",
		Status:  http.StatusInternalServerError,
	}

	ErrInvalidJSON = &AppError{
		Code:    "INVALID_JSON",
		Message: "Invalid JSON body format",
		Status:  http.StatusBadRequest,
	}
)
