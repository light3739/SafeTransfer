// internal/model/file.go

package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	CID           string `gorm:"column:cid;type:varchar(255);unique_index"`
	EncryptionKey string `gorm:"type:varchar(255)"`
	Nonce         string `gorm:"type:varchar(255)"`
	Signature     string `gorm:"type:text"`
}
