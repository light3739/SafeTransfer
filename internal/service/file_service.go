package service

import (
	"SafeTransfer/internal/crypto"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
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
	FileRepo    repository.FileRepository // Use FileRepository instead of direct database access
	PrivateKey  *rsa.PrivateKey
}

// NewFileService creates a new instance of FileService with dependencies injected.
func NewFileService(ipfsStorage *storage.IPFSStorage, fileRepo repository.FileRepository, privateKey *rsa.PrivateKey) *FileService {
	return &FileService{
		IPFSStorage: ipfsStorage,
		FileRepo:    fileRepo,
		PrivateKey:  privateKey,
	}
}

// UploadFile handles the uploading of a file, including processing, encryption, and storage.
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
	fileMetadata := &model.File{ // Note: Use a pointer to match the repository interface
		CID:           cid,
		EncryptionKey: base64.StdEncoding.EncodeToString(key),
		Nonce:         nonceStr,
		Signature:     signatureStr,
	}

	// Use the FileRepository to save file metadata
	if err := fs.FileRepo.SaveFileMetadata(fileMetadata); err != nil {
		return "", err
	}

	return cid, nil
}

// processFile handles the processing of the file, including signing and encryption key generation.
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
