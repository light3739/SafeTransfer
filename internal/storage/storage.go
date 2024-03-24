// internal/storage/storage.go

package storage

import (
	"SafeTransfer/internal/crypto" // Import the crypto package
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
	ipfsShell := shell.NewShell(apiURL)
	return &IPFSStorage{shell: ipfsShell}
}

// UploadFileToIPFS uploads a file to IPFS and returns the generated CID and the nonce.
func (is *IPFSStorage) UploadFileToIPFS(file io.ReadSeeker, key []byte) (string, []byte, error) {
	// Encrypt the file before uploading
	encryptedFile, nonce, err := crypto.EncryptFile(file, key)
	if err != nil {
		return "", nil, fmt.Errorf("failed to encrypt file: %w", err)
	}
	defer func() {
		if closer, ok := encryptedFile.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				fmt.Println("Error closing encryptedFile:", err)
			}
		}
	}()

	// Reset the file reader to the beginning
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", nil, fmt.Errorf("failed to reset file reader: %w", err)
	}

	// Add the encrypted file to IPFS
	cid, addErr := is.shell.Add(encryptedFile)
	if addErr != nil {
		return "", nil, fmt.Errorf("failed to upload encrypted file to IPFS: %w", addErr)
	}

	return cid, nonce, nil
}

// DownloadFileFromIPFS retrieves a file from IPFS using its CID.
func (is *IPFSStorage) DownloadFileFromIPFS(cid string) (io.ReadCloser, error) {
	// Use the IPFS shell to retrieve the encrypted file
	reader, err := is.shell.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("failed to download encrypted file from IPFS: %w", err)
	}
	readCloser := io.NopCloser(reader)

	return readCloser, nil
}
