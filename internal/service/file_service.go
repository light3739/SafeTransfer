package service

import (
	"SafeTransfer/internal/crypto"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
	"SafeTransfer/internal/storage"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
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
	FileRepo    repository.FileRepository
}

// NewFileService creates a new instance of FileService with dependencies injected.
func NewFileService(ipfsStorage *storage.IPFSStorage, fileRepo repository.FileRepository) *FileService {
	return &FileService{
		IPFSStorage: ipfsStorage,
		FileRepo:    fileRepo,
	}
}

// UploadFile handles the uploading of a file, including processing, encryption, and storage.
func (fs *FileService) UploadFile(file multipart.File, ethereumAddress string) (string, string, error) {
	// Generate a new key pair for each file
	privateKey, err := generateRSAKeyPair(2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	signatureStr, publicKeyStr, key, originalFileHash, err := fs.processFile(file, privateKey)
	if err != nil {
		return "", "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	// Verify the file signature
	publicKey, err := parsePublicKey(publicKeyStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse public key: %w", err)
	}
	if err := crypto.VerifyFile(file, signatureStr, publicKey); err != nil {
		return "", "", fmt.Errorf("file verification failed: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", fmt.Errorf("failed to reset file reader: %w", err)
	}

	cid, nonce, err := fs.IPFSStorage.UploadFileToIPFS(file, key)
	if err != nil {
		return "", "", err
	}

	nonceStr := base64.StdEncoding.EncodeToString(nonce)
	fileMetadata := &model.File{
		CID:             cid,
		EncryptionKey:   base64.StdEncoding.EncodeToString(key),
		Nonce:           nonceStr,
		Signature:       signatureStr,
		EthereumAddress: ethereumAddress,
		PublicKey:       publicKeyStr,
	}

	if err := fs.FileRepo.SaveFileMetadata(fileMetadata); err != nil {
		return "", "", err
	}

	return cid, base64.StdEncoding.EncodeToString(originalFileHash), nil
}

// parsePublicKey parses a base64-encoded public key string into an *rsa.PublicKey.
func parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an RSA public key")
	}

	return rsaPublicKey, nil
}

// processFile handles the processing of the file, including signing and encryption key generation.
func (fs *FileService) processFile(file io.ReadSeeker, privateKey *rsa.PrivateKey) (signatureStr string, publicKeyStr string, key []byte, originalFileHash []byte, err error) {
	// Hash the file content
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to hash file: %w", err)
	}
	originalFileHash = hash.Sum(nil)

	// Reset the file reader to the beginning for signing
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to reset file reader: %w", err)
	}

	// Sign the file
	signature, err := crypto.SignFile(file, privateKey)
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to sign file: %w", err)
	}
	signatureStr = base64.StdEncoding.EncodeToString(signature)

	// Encode the public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyStr = base64.StdEncoding.EncodeToString(publicKeyBytes)

	// Generate an encryption key
	key = make([]byte, encryptionKeySize)
	if _, err := rand.Read(key); err != nil {
		return "", "", nil, nil, fmt.Errorf("failed to generate encryption key: %w", err)
	}

	return signatureStr, publicKeyStr, key, originalFileHash, nil
}

// generateRSAKeyPair generates a new RSA key pair with the specified key size.
func generateRSAKeyPair(keySize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}
	return privateKey, nil
}
