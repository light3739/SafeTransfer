package service

import (
	"SafeTransfer/internal/repository"
	"SafeTransfer/internal/storage"
	"encoding/base64"
	"fmt"
	"io"
)

type DownloadService struct {
	IPFSStorage *storage.IPFSStorage
	FileRepo    repository.FileRepository // Use FileRepository instead of direct database access
}

// NewDownloadService creates a new instance of DownloadService with dependencies injected.
func NewDownloadService(ipfsStorage *storage.IPFSStorage, fileRepo repository.FileRepository) *DownloadService {
	return &DownloadService{
		IPFSStorage: ipfsStorage,
		FileRepo:    fileRepo,
	}
}

// DownloadFile handles the downloading of a file by its CID.
func (ds *DownloadService) DownloadFile(cid string) (io.ReadCloser, error) {
	// Retrieve the stored file metadata from the repository
	fileMetadata, err := ds.FileRepo.GetFileMetadataByCID(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file metadata: %w", err)
	}

	// Decode the stored encryption key and nonce
	encryptionKey, err := base64.StdEncoding.DecodeString(fileMetadata.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encryption key: %w", err)
	}
	nonce, err := base64.StdEncoding.DecodeString(fileMetadata.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %w", err)
	}

	// Download the file from IPFS using the key and nonce
	reader, err := ds.IPFSStorage.DownloadFileFromIPFS(cid, encryptionKey, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return io.NopCloser(reader), nil
}
