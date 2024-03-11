package api

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/storage"
	"crypto/rsa"
	"github.com/go-chi/chi"
	"net/http"
)

// Handler represents the handler for API endpoints.
type Handler struct {
	ipfsStorage *storage.IPFSStorage
	db          *db.Database
	privateKey  *rsa.PrivateKey
}

// NewAPIHandler creates a new instance of APIHandler.
func NewAPIHandler(ipfsStorage *storage.IPFSStorage, db *db.Database, privateKey *rsa.PrivateKey) *Handler {
	return &Handler{
		ipfsStorage: ipfsStorage,
		db:          db,
		privateKey:  privateKey,
	}
}

// RegisterRoutes registers API routes to the provided router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/upload", h.HandleFileUpload)
	r.Get("/download/{cid}", h.HandleFileDownload)
}

// HandleFileUpload handles the file upload endpoint.
func (h *Handler) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	HandleFileUpload(w, r, h.ipfsStorage, h.privateKey, h.db)
}
