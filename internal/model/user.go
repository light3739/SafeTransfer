// internal/model/user.go

package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	EthereumAddress string `gorm:"uniqueIndex"` // Unique Ethereum address of the user
	Nonce           string // Nonce for authentication
}
