package handlers

import (
	"encoding/json"
	"errors"
	"kogatari/internal/internal"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		var msg string
		var status int
		var userExists bool
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			msg = "user with this email already exist"
			status = http.StatusConflict
			userExists = true
		} else {
			msg = "user with this email already exist"
			status = http.StatusInternalServerError
		}
		app.Slogger.ErrorContext(r.Context(),
			"error with creating user has occured",
			slog.String("err", err.Error()),
			slog.Bool("user exists in DB", userExists),
			slog.String("email", params.Email))

		internal.RespondWithError(w, r, status, msg)
		return
	}

	internal.RespondWithJSON(w, r, http.StatusCreated, usr)
}
