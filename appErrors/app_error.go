package appErrors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
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
)
