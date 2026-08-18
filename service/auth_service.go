package service

import (
	"net/http"
	"url-shortener/appErrors"
	"url-shortener/config"
	"url-shortener/dto"
	"url-shortener/models"
	"url-shortener/repository"
	"url-shortener/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	ValidateToken(token string) bool
	GetUserFromToken(token string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (a authService) ValidateToken(token string) bool {
	//TODO implement me
	panic("implement me")
}

func (a authService) GetUserFromToken(token string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (a authService) Register(req dto.RegisterRequest) error {
	if _, err := a.userRepo.FindByEmail(req.Email); err == nil {
		return &appErrors.AppError{
			Code:    "EMAIL_ALREADY_EXISTS",
			Message: "User already exists with this email",
			Status:  http.StatusConflict,
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return appErrors.ErrServerError
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := a.userRepo.Create(user); err != nil {
		return appErrors.ErrServerError
	}

	return nil
}

func (a authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, appErrors.ErrUserNotFound
	}

	//Comparing passwords
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, &appErrors.AppError{
			Code:    "WRONG_PASSWORD",
			Message: "Wrong password",
			Status:  http.StatusBadRequest,
		}
	}

	token, err := utils.GenerateToken(user.ID, user.Email, a.cfg.JwtSecret, 1)
	if err != nil {
		return nil, &appErrors.AppError{
			Code:    "FAILED_TOKEN_GENERATION",
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}

	return &dto.AuthResponse{
		Token: token,
		User:  *user,
	}, nil

}
