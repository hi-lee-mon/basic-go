package services

import (
	"basic-go/src/models"
	"basic-go/src/repositories"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func createToken(userId uint, email string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId, // JWTの標準クレームで、ユーザーIDを指定(subjectの略)
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(), // 1時間後にトークンが期限切れになるように設定(expはexpirationの略)
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

type IAuthService interface {
	Signup(email, password string) error
	Login(email, password string) (*string, error)
	GetUserFromToken(tokenString string) (*models.User, error)
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

func (s *AuthService) Login(email, password string) (*string, error) {
	foundUser, err := s.repository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		// 404
		return nil, errors.New("invalid email or password")
	}
	return createToken(foundUser.ID, foundUser.Email)
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 1.型アサーションをする。SigningMethodHMACはHMAC署名アルゴリズムを使用していることを確認するためのもの。
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	var user *models.User

	// 2.型アサーションをして、トークンが有効であることを確認する。MapClaimsはJWTのクレームをマップ形式で表現するためのもの。
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, jwt.ErrTokenExpired
		}

		// 3.jwtからユーザ情報を取得して、DBからユーザを取得する
		user, err = s.repository.FindUserByEmail(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
