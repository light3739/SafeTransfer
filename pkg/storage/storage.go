// Package storage provides a service for interacting with IPFS for file storage.
package storage

import (
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
func (is *IPFSStorage) UploadFileToIPFS(file io.Reader) (string, error) {
	// Read the file data
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Log the file data length for debugging
	fmt.Printf("File data length: %d bytes\n", len(data))

	// Add the file to IPFS
	fmt.Println("Uploading file to IPFS...")
	cid, err := is.shell.Add(bytes.NewReader(data))
	if err != nil {
		fmt.Printf("Failed to upload file to IPFS: %v\n", err)
		return "", fmt.Errorf("failed to upload file to IPFS: %w", err)
	}

	fmt.Printf("File uploaded successfully. CID: %s\n", cid)

	return cid, nil
}
