// response/response.go

package api

import (
	"encoding/json"
	"io"
	"net/http"
)

// RespondWithError sends an error response in JSON format.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response with the provided data.
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

// SendFile sends the file data in the response, including the SHA-256 hash of the file content in the headers.
func SendFile(w http.ResponseWriter, reader io.Reader, filename string, hash string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-File-Hash", hash) // Include the SHA-256 hash in the response headers
	_, err := io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
		return
	}
}
