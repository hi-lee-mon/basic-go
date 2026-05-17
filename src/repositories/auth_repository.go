package repositories

import (
	"basic-go/src/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(user models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user models.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	// 第2引数はwhere句の条件を指定するためのもので、?はプレースホルダで第3引数のemailの値が?に置き換えられる。
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		return nil, result.Error
	}
	return &user, nil
}
