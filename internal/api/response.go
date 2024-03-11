// response/response.go

package api

import (
	"encoding/json"
	"io"
	"net/http"
)

// RespondWithError sends an error response in JSON format.
func RespondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response with the provided data.
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendFile sends the file data in the response.
func SendFile(w http.ResponseWriter, reader io.Reader, filename string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Write the file data to the response
	io.Copy(w, reader)
}
