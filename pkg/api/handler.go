// pkg/api/handler.go

package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"mime/multipart"
	"net/http"

	"SafeTransfer/pkg/storage"
)

// APIHandler represents the handler for API endpoints.
type APIHandler struct {
	ipfsStorage *storage.IPFSStorage
}

// NewAPIHandler creates a new instance of APIHandler.
func NewAPIHandler(ipfsStorage *storage.IPFSStorage) *APIHandler {
	return &APIHandler{ipfsStorage: ipfsStorage}
}

// FileUploadRequest represents the request structure for file upload.
type FileUploadRequest struct {
	File *multipart.FileHeader `json:"file"`
}

// FileUploadResponse represents the response structure for file upload.
type FileUploadResponse struct {
	CID string `json:"cid"`
}

// RegisterRoutes registers API routes to the provided router.
func (ah *APIHandler) RegisterRoutes(r chi.Router) {
	r.HandleFunc("/upload", ah.HandleFileUpload)
}

// HandleFileUpload handles the file upload endpoint.
// @Summary Uploads a file to IPFS
// @Description Uploads a file to IPFS and returns the CID.
// @Tags File
// @Accept application/json
// @Produce application/json
// @Param file formData file true "File to upload"
// @Success   200 {object} FileUploadResponse
// @Failure   400 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router /upload [post]
func (ah *APIHandler) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for the file
	if err != nil {
		ah.respondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		ah.respondWithError(w, http.StatusBadRequest, "Failed to get file from form data")
		return
	}
	defer file.Close()

	// Upload file to IPFS
	cid, err := ah.ipfsStorage.UploadFileToIPFS(file)
	if err != nil {
		ah.respondWithError(w, http.StatusInternalServerError, "Failed to upload file to IPFS")
		return
	}

	// Respond with CID
	response := FileUploadResponse{CID: cid}
	ah.respondWithJSON(w, http.StatusOK, response)
}

// respondWithError sends an error response in JSON format.
func (ah *APIHandler) respondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondWithJSON sends a JSON response with the provided data.
func (ah *APIHandler) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
