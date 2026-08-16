package service

import (
	"errors"
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
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
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
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}

	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if _, err := a.userRepo.Create(user); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (a authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := a.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	//Comparing passwords
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	//TODO: Implement jwt token generation
	token := "jwt_token"

	return &dto.AuthResponse{
		Token: token,
		User:  *user,
	}, nil

}
