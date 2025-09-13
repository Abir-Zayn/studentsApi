package response

import (
	"encoding/json"
	"net/http"
)

// Response represents the standard API response structure
type Response struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"status_code"`
}

// Status constants
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

// WriteJson writes a JSON response with proper headers and status code
func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// Send sends a response with appropriate status code
func Send(w http.ResponseWriter, resp Response) error {
	return WriteJson(w, resp.StatusCode, resp)
}
