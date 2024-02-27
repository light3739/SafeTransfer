// upload/upload.go

package upload

import (
	"SafeTransfer/pkg/api/response"
	"SafeTransfer/pkg/crypto"
	"SafeTransfer/pkg/db"
	"SafeTransfer/pkg/model"
	"SafeTransfer/pkg/storage"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
)

// FileUploadRequest represents the request for file upload.
type FileUploadRequest struct {
	File *multipart.FileHeader `json:"file"`
}

// FileUploadResponse represents the response for file upload.
type FileUploadResponse struct {
	CID string `json:"cid"`
}

// HandleFileUpload handles the file upload endpoint.
// @Summary Uploads a file to IPFS
// @Description Uploads a file to IPFS and returns the CID.
// @Tags File
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "File to upload"
// @Success   200 {object} FileUploadResponse
// @Failure   400 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router /upload [post]
// HandleFileUpload handles the file upload endpoint.
func HandleFileUpload(w http.ResponseWriter, r *http.Request, ipfsStorage *storage.IPFSStorage, privateKey *rsa.PrivateKey, db *db.Database) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Failed to get file from form data")
		return
	}
	defer file.Close()

	// Read the entire file into a buffer
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Create a new bytes.Reader for each operation that needs to read the file
	fileReader := bytes.NewReader(fileBytes)

	// Sign the file before uploading
	signature, err := crypto.SignFile(fileReader, privateKey)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to sign file")
		return
	}
	signatureStr := base64.StdEncoding.EncodeToString([]byte(signature))

	// Generate an encryption key
	key := make([]byte, 32) // Example:   32 bytes for AES-256
	if _, err := rand.Read(key); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to generate encryption key")
		return
	}

	// Create a new bytes.Reader for the upload operation
	// This ensures that the file is read from the beginning
	uploadFileReader := bytes.NewReader(fileBytes)

	// Upload the file to IPFS
	cid, nonce, err := ipfsStorage.UploadFileToIPFS(uploadFileReader, key)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to upload file to IPFS")
		return
	}

	// Convert nonce to a string for storage
	nonceStr := base64.StdEncoding.EncodeToString(nonce)

	// Save file metadata to the database
	fileMetadata := model.File{
		CID:           cid,
		EncryptionKey: base64.StdEncoding.EncodeToString(key),
		Nonce:         nonceStr,
		Signature:     signatureStr,
	}
	if err := db.SaveFileMetadata(fileMetadata); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to save file metadata")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, FileUploadResponse{CID: cid})
}
