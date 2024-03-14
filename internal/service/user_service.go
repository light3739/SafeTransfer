// internal/service/user_service.go

package service

import (
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
	"crypto/rand"
	"encoding/hex"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

// GenerateNonceForUser generates a unique nonce for the user and stores it in the database.
func (us *UserService) GenerateNonceForUser(ethereumAddress string) (string, error) {
	nonceBytes := make([]byte, 16) // 128-bit nonce
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(nonceBytes)

	// Store or update the nonce for the user
	user := &model.User{
		EthereumAddress: ethereumAddress,
		Nonce:           nonce,
	}
	err := us.UserRepo.SaveOrUpdateUser(user)
	return nonce, err
}
