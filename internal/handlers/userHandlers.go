package handlers

import (
	"encoding/json"
	"kogatari/internal/internal"
	"log/slog"
	"net/http"
)

func (app application) CreateUserHander(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Email string `json:"email"`
	}{}

	err := decoder.Decode(&params)
	if err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"error with request body has occured",
			slog.String("err", err.Error()))

		internal.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	}

	usr, err := app.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"error with creating user has occured",
			slog.String("err", err.Error()),
			slog.String("email", params.Email))

		internal.RespondWithError(w, r, http.StatusInternalServerError, "internal error")
		return
	}

	internal.RespondWithJSON(w, r, http.StatusCreated, usr)
}
