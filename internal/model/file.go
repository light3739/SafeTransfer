// internal/model/file.go

package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	CID             string `gorm:"column:cid;type:varchar(255);uniqueIndex"`
	EthereumAddress string `gorm:"column:ethereum_address;type:varchar(255);index"`
	EncryptionKey   string `gorm:"column:encryption_key;type:varchar(255)"`
	Nonce           string `gorm:"column:nonce;type:varchar(255)"`
	Signature       string `gorm:"column:signature;type:text"`
}
