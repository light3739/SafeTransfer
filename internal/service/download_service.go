package service

import (
	"SafeTransfer/internal/crypto"
	"SafeTransfer/internal/repository"
	"SafeTransfer/internal/storage"
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
)

type DownloadService struct {
	IPFSStorage *storage.IPFSStorage
	FileRepo    repository.FileRepository
}

// NewDownloadService creates a new instance of DownloadService with dependencies injected.
func NewDownloadService(ipfsStorage *storage.IPFSStorage, fileRepo repository.FileRepository) *DownloadService {
	return &DownloadService{
		IPFSStorage: ipfsStorage,
		FileRepo:    fileRepo,
	}
}

// DownloadFile handles the downloading of a file by its CID.
func (ds *DownloadService) DownloadFile(cid string) (io.Reader, error) {
	fileMetadata, err := ds.FileRepo.GetFileMetadataByCID(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(fileMetadata.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(fileMetadata.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	encryptedFile, err := ds.IPFSStorage.DownloadFileFromIPFS(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from IPFS: %w", err)
	}
	defer encryptedFile.Close()

	// First, decrypt the file
	decryptedContent, err := crypto.DecryptFile(encryptedFile, encryptionKey, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %w", err)
	}

	decryptedData, err := io.ReadAll(decryptedContent)
	if err != nil {
		return nil, fmt.Errorf("failed to read decrypted file content: %w", err)
	}

	// Then, verify the signature
	publicKeyBytes, err := base64.StdEncoding.DecodeString(fileMetadata.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}

	// Convert decrypted data back to a reader for signature verification
	decryptedReader := bytes.NewReader(decryptedData)
	if err := crypto.VerifyFile(decryptedReader, fileMetadata.Signature, rsaPublicKey); err != nil {
		return nil, fmt.Errorf("file verification failed: %w", err)
	}

	// Return the decrypted content as a reader
	return bytes.NewReader(decryptedData), nil
}
