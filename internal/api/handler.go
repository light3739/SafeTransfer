package api

import (
	"SafeTransfer/internal/service"
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strings"
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
	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware)

		r.Post("/upload", h.handleFileUpload)
		r.Get("/download/{cid}", h.handleFileDownload)
		r.Get("/checkToken", h.handleCheckToken)
	})
	r.Post("/verifySignature", h.handleVerifySignature)
	r.Post("/generateNonce", h.handleGenerateNonce)
}
func (h *Handler) handleCheckToken(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "This is a test message for authenticated users."})
}

// @Summary Upload a file
// @Description Uploads a file and associates it with the user's Ethereum address
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /upload [post]
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

	ethereumAddress := r.Header.Get("EthereumAddress")
	cid, originalFileHash, err := h.FileService.UploadFile(file, ethereumAddress)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{
		"cid":              cid,
		"originalFileHash": originalFileHash,
	}

	RespondWithJSON(w, http.StatusOK, response)
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

	SendFile(w, reader, cid)
}

func (h *Handler) handleVerifySignature(w http.ResponseWriter, r *http.Request) {
	type request struct {
		EthereumAddress string `json:"ethereumAddress"`
		Signature       string `json:"signature"`
		Message         string `json:"message"` // Add the message field
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Received address: %s, signature: %s, message: %s\n", req.EthereumAddress, req.Signature, req.Message) // Update the log to include the message

	// Verify the signature against the message instead of the nonce
	recoveredAddress, err := h.UserService.VerifySignature(req.Message, req.Signature)
	if err != nil || strings.ToLower(recoveredAddress) != strings.ToLower(req.EthereumAddress) {
		RespondWithError(w, http.StatusUnauthorized, "Invalid signature")
		return
	}

	token, err := h.UserService.GenerateJWT(req.EthereumAddress)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate JWT token")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleGenerateNonce(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EthereumAddress string `json:"ethereumAddress"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	nonce, err := h.UserService.GenerateNonceForUser(req.EthereumAddress)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate nonce")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"nonce": nonce})
}
