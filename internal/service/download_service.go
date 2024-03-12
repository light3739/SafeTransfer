package service

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/storage"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil" // Import ioutil for NopCloser
)

type DownloadService struct {
	IPFSStorage *storage.IPFSStorage
	DB          *db.Database
}

func NewDownloadService(ipfsStorage *storage.IPFSStorage, db *db.Database) *DownloadService {
	return &DownloadService{
		IPFSStorage: ipfsStorage,
		DB:          db,
	}
}

func (ds *DownloadService) DownloadFile(cid string) (io.ReadCloser, error) {
	// Retrieve the stored file metadata from the database
	fileMetadata, err := ds.DB.GetFileMetadataByCID(cid)
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

	// Wrap the reader with NopCloser to satisfy the io.ReadCloser interface
	return ioutil.NopCloser(reader), nil
}
