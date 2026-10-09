package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("respondWithJSON: marshal failed", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(`{"error":"internal server error"}`)); err != nil {
			slog.ErrorContext(r.Context(), "respondWithJSON: write failed", slog.String("error", err.Error()))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		slog.Error("respondWithJSON: write failed", "error", err)
	}
}

func RespondWithError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	RespondWithJSON(w, r, code, errorResponse{Error: msg})
}
