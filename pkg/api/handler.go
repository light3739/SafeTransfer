package api

import (
	"SafeTransfer/pkg/api/download"
	"SafeTransfer/pkg/api/upload"
	"SafeTransfer/pkg/db"
	"SafeTransfer/pkg/storage"
	"crypto/rsa"
	"github.com/go-chi/chi"
	"net/http"
)

// Handler represents the handler for API endpoints.
type Handler struct {
	ipfsStorage     *storage.IPFSStorage
	downloadHandler *download.Handler
	db              *db.Database
	privateKey      *rsa.PrivateKey // Add this line
}

// NewAPIHandler creates a new instance of APIHandler.
func NewAPIHandler(ipfsStorage *storage.IPFSStorage, db *db.Database, privateKey *rsa.PrivateKey) *Handler {
	downloadHandler := download.NewDownloadHandler(ipfsStorage, db)
	return &Handler{
		ipfsStorage:     ipfsStorage,
		downloadHandler: downloadHandler,
		db:              db,
		privateKey:      privateKey,
	}
}

// RegisterRoutes registers API routes to the provided router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/upload", h.HandleFileUpload)
	r.Get("/download/{cid}", h.HandleFileDownload)
}

// HandleFileUpload handles the file upload endpoint.
func (h *Handler) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	upload.HandleFileUpload(w, r, h.ipfsStorage, h.privateKey, h.db)
}

// HandleFileDownload handles the file download process.
func (h *Handler) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	h.downloadHandler.HandleFileDownload(w, r)
}
