package service

import (
	"SafeTransfer/internal/crypto"
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/storage"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
)

const (
	MaxMultipartFormSize = 10 << 20 // 10 MB
	encryptionKeySize    = 32       // 32 bytes for AES-256
)

type FileService struct {
	IPFSStorage *storage.IPFSStorage
	DB          *db.Database
	PrivateKey  *rsa.PrivateKey
}

func NewFileService(ipfsStorage *storage.IPFSStorage, db *db.Database, privateKey *rsa.PrivateKey) *FileService {
	return &FileService{
		IPFSStorage: ipfsStorage,
		DB:          db,
		PrivateKey:  privateKey,
	}
}

func (fs *FileService) UploadFile(file multipart.File) (string, error) {
	signatureStr, key, err := fs.processFile(file)
	if err != nil {
		return "", err
	}

	// Reset file reader to the beginning for upload
	file.Seek(0, io.SeekStart)

	cid, nonce, err := fs.IPFSStorage.UploadFileToIPFS(file, key)
	if err != nil {
		return "", err
	}

	nonceStr := base64.StdEncoding.EncodeToString(nonce)
	fileMetadata := model.File{
		CID:           cid,
		EncryptionKey: base64.StdEncoding.EncodeToString(key),
		Nonce:         nonceStr,
		Signature:     signatureStr,
	}

	if err := fs.DB.SaveFileMetadata(fileMetadata); err != nil {
		return "", err
	}

	return cid, nil
}

func (fs *FileService) processFile(file multipart.File) (signatureStr string, key []byte, err error) {
	// Sign the file
	signature, err := crypto.SignFile(file, fs.PrivateKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign file: %w", err)
	}
	signatureStr = base64.StdEncoding.EncodeToString([]byte(signature))

	// Generate an encryption key
	key = make([]byte, encryptionKeySize)
	if _, err := rand.Read(key); err != nil {
		return "", nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	return signatureStr, key, nil
}
