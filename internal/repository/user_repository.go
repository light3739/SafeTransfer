package repository

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveOrUpdateUser(user *model.User) error
	FindByEthereumAddress(ethereumAddress string) (*model.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *db.Database) UserRepository {
	return &UserRepositoryImpl{DB: db.DB}
}

func (repo *UserRepositoryImpl) SaveOrUpdateUser(user *model.User) error {
	var existingUser model.User
	result := repo.DB.Where("ethereum_address = ?", user.EthereumAddress).First(&existingUser)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.DB.Create(user).Error
	} else if result.Error != nil {
		return result.Error
	}

	existingUser.Nonce = user.Nonce
	return repo.DB.Save(&existingUser).Error
}

func (repo *UserRepositoryImpl) FindByEthereumAddress(ethereumAddress string) (*model.User, error) {
	var user model.User
	result := repo.DB.Where("ethereum_address = ?", ethereumAddress).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
