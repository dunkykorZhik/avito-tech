package routes

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func errorJSON(w http.ResponseWriter, msg string, status int) {
	type envelope struct {
		ErrorMsg string `json:"errors"`
	}

	if err := writeJSON(w, &envelope{ErrorMsg: msg}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)

}
