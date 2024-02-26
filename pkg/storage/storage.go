// pkg/storage/storage.go

package storage

import (
	"SafeTransfer/pkg/crypto" // Import the crypto package
	"bytes"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
	"io"
)

// IPFSStorage represents the IPFS storage service.
type IPFSStorage struct {
	shell *shell.Shell
}

// NewIPFSStorage creates a new instance of IPFSStorage.
func NewIPFSStorage(apiURL string) *IPFSStorage {
	shell := shell.NewShell(apiURL)
	return &IPFSStorage{shell: shell}
}

// UploadFileToIPFS uploads a file to IPFS and returns the generated CID.
func (is *IPFSStorage) UploadFileToIPFS(file io.Reader, key []byte) (string, error) {
	// Encrypt the file before uploading
	encryptedFile, err := crypto.EncryptFile(file, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt file: %w", err)
	}

	// Read the encrypted file data
	data, err := io.ReadAll(encryptedFile)
	if err != nil {
		return "", fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Add the encrypted file to IPFS
	cid, err := is.shell.Add(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to upload encrypted file to IPFS: %w", err)
	}

	return cid, nil
}

// DownloadFileFromIPFS retrieves a file from IPFS using its CID.
func (is *IPFSStorage) DownloadFileFromIPFS(cid string, key []byte, nonce []byte) (io.Reader, error) {
	// Use the IPFS shell to retrieve the encrypted file
	reader, err := is.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to download encrypted file from IPFS: %w", err)
	}

	// Decrypt the file after downloading
	decryptedFile, err := crypto.DecryptFile(reader, key, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt file: %w", err)
	}

	return decryptedFile, nil
}
