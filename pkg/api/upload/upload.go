package upload

import (
	"SafeTransfer/pkg/api/response"
	"SafeTransfer/pkg/crypto"
	"SafeTransfer/pkg/db"
	"SafeTransfer/pkg/model"
	"SafeTransfer/pkg/storage"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	maxMultipartFormSize = 10 << 20 // 10 MB
	encryptionKeySize    = 32       // 32 bytes for AES-256
)

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
	if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Failed to get file from form data")
		return
	}
	defer file.Close()

	signatureStr, key, err := processFile(file, privateKey)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Reset file reader to the beginning for upload
	file.Seek(0, io.SeekStart)

	cid, nonce, err := ipfsStorage.UploadFileToIPFS(file, key)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to upload file to IPFS")
		return
	}

	nonceStr := base64.StdEncoding.EncodeToString(nonce)
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

func processFile(file multipart.File, privateKey *rsa.PrivateKey) (signatureStr string, key []byte, err error) {
	// Sign the file
	signature, err := crypto.SignFile(file, privateKey)
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
