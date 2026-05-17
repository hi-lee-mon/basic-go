package services

import (
	"basic-go/src/models"
	"basic-go/src/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email, password string) error
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(authRepo repositories.IAuthRepository) *AuthService {
	return &AuthService{repository: authRepo}
}

func (s *AuthService) Signup(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repository.CreateUser(user)
}
