package helpers

import (
	"encoding/json"
	"net/http"
)

// error response struct
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RespondWithError writes a JSON-encoded error response
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
