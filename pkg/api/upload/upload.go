// upload/upload.go

package upload

import (
	"SafeTransfer/pkg/api/response"
	"SafeTransfer/pkg/storage"
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
// @Success  200 {object} FileUploadResponse
// @Failure  400 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router /upload [post]
func HandleFileUpload(w http.ResponseWriter, r *http.Request, ipfsStorage *storage.IPFSStorage) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.RespondWithError(w, http.StatusBadRequest, "Failed to get file from form data")
		return
	}
	defer file.Close()

	cid, err := ipfsStorage.UploadFileToIPFS(file)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to upload file to IPFS")
		return
	}

	response.RespondWithJSON(w, http.StatusOK, FileUploadResponse{CID: cid})
}
