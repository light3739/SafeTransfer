// download/download.go

package api

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/storage"
	"encoding/base64"
	"github.com/go-chi/chi"
	"net/http"
)

// NewDownloadHandler creates a new instance of DownloadHandler.
func NewDownloadHandler(ipfsStorage *storage.IPFSStorage, db *db.Database) *Handler {
	return &Handler{ipfsStorage: ipfsStorage, db: db} // Pass the db instance here
}

// HandleFileDownload handles the file download process.
// @Summary Download a file
// @Description Downloads a file from IPFS using its CID
// @ID download-file
// @Accept  json
// @Produce  application/octet-stream
// @Param cid path string true "CID of the file to download"
// @Success  200 {file} file "File data"
// @Failure  400 {object} string "Bad Request"
// @Failure  500 {object} string "Internal Server Error"
// @Router /download/{cid} [get]
func (h *Handler) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		RespondWithError(w, http.StatusBadRequest, "CID is required")
		return
	}

	// Retrieve the stored file metadata from the database
	fileMetadata, err := h.db.GetFileMetadataByCID(cid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve file metadata")
		return
	}

	// Decode the stored encryption key and nonce
	encryptionKey, err := base64.StdEncoding.DecodeString(fileMetadata.EncryptionKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to decode encryption key")
		return
	}
	nonce, err := base64.StdEncoding.DecodeString(fileMetadata.Nonce)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to decode nonce")
		return
	}

	// Download the file from IPFS using the key and nonce
	reader, err := h.ipfsStorage.DownloadFileFromIPFS(cid, encryptionKey, nonce)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to download file")
		return
	}

	// Assuming you no longer need to verify the file's signature during download
	// Send the file to the client
	SendFile(w, reader, cid)
}
