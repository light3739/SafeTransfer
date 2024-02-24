// api/handler.go

package api

import (
	"SafeTransfer/pkg/api/download"
	"SafeTransfer/pkg/api/upload"
	"SafeTransfer/pkg/storage"
	"github.com/go-chi/chi"
	"net/http"
)

// Handler represents the handler for API endpoints.
type Handler struct {
	ipfsStorage     *storage.IPFSStorage
	downloadHandler *download.Handler
}

// NewAPIHandler creates a new instance of APIHandler.
func NewAPIHandler(ipfsStorage *storage.IPFSStorage) *Handler {
	return &Handler{
		ipfsStorage:     ipfsStorage,
		downloadHandler: download.NewDownloadHandler(ipfsStorage),
	}
}

// RegisterRoutes registers API routes to the provided router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.HandleFunc("/upload", h.HandleFileUpload)
	r.HandleFunc("/download/{cid}", h.HandleFileDownload)
}

func (h *Handler) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	upload.HandleFileUpload(w, r, h.ipfsStorage)
}

func (h *Handler) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	h.downloadHandler.HandleFileDownload(w, r)
}
