package api

import (
	"SafeTransfer/internal/service"
	"github.com/go-chi/chi"
	"net/http"
)

type Handler struct {
	FileService     *service.FileService
	DownloadService *service.DownloadService
	UserService     *service.UserService
}

func NewAPIHandler(fileService *service.FileService, downloadService *service.DownloadService, userService *service.UserService) *Handler {
	return &Handler{
		FileService:     fileService,
		DownloadService: downloadService,
		UserService:     userService,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/upload", h.handleFileUpload)
	r.Get("/download/{cid}", h.handleFileDownload)
	r.Get("/nonce/{ethereumAddress}", h.handleGetNonce)
}

func (h *Handler) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(service.MaxMultipartFormSize); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to parse form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Failed to get file from form data")
		return
	}
	defer file.Close()

	cid, err := h.FileService.UploadFile(file)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"cid": cid})
}

func (h *Handler) handleFileDownload(w http.ResponseWriter, r *http.Request) {
	cid := chi.URLParam(r, "cid")
	if cid == "" {
		RespondWithError(w, http.StatusBadRequest, "CID is required")
		return
	}

	reader, err := h.DownloadService.DownloadFile(cid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	SendFile(w, reader, cid)
}

func (h *Handler) handleGetNonce(w http.ResponseWriter, r *http.Request) {
	ethereumAddress := chi.URLParam(r, "ethereumAddress")
	if ethereumAddress == "" {
		RespondWithError(w, http.StatusBadRequest, "Ethereum address is required")
		return
	}

	nonce, err := h.UserService.GenerateNonceForUser(ethereumAddress)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate nonce")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"nonce": nonce})
}
