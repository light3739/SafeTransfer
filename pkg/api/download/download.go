// download/download.go

package download

import (
	"SafeTransfer/pkg/api/response"
	"SafeTransfer/pkg/storage"
	"github.com/go-chi/chi"
	"net/http"
)

// Handler represents the handler for file download endpoints.
type Handler struct {
	ipfsStorage *storage.IPFSStorage
}

// NewDownloadHandler creates a new instance of DownloadHandler.
func NewDownloadHandler(ipfsStorage *storage.IPFSStorage) *Handler {
	return &Handler{ipfsStorage: ipfsStorage}
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
		response.RespondWithError(w, http.StatusBadRequest, "CID is required")
		return
	}

	reader, err := h.ipfsStorage.DownloadFileFromIPFS(cid)
	if err != nil {
		response.RespondWithError(w, http.StatusInternalServerError, "Failed to download file")
		return
	}

	response.SendFile(w, reader, cid)
}
