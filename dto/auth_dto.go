package dto

import "url-shortener/models"

type RegisterRequest struct {
	Email           string `json:"email" form:"email" validate:"required,email"`
	Password        string `json:"password" form:"password" validate:"required,min=8,password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
