package routes

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	return decoder.Decode(data)
}

type ErrorResponse struct {
	ErrorMsg string `json:"errors"`
}

func errorJSON(w http.ResponseWriter, msg string, status int) {

	if err := writeJSON(w, status, &ErrorResponse{ErrorMsg: msg}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
