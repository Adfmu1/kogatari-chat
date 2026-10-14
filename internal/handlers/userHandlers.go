package handlers

import (
	"encoding/json"
	"errors"
	"kogatari/internal/auth"
	"kogatari/internal/database"
	"kogatari/internal/httputil"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func (app application) CreateUserHander(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params database.CreateUserParams
	var err error

	if err := decoder.Decode(&params); err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"error with request body has occured",
			slog.String("err", err.Error()))

		httputil.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	}

	if passLength := len(params.HashedPassword); passLength < 8 {
		app.Slogger.InfoContext(r.Context(),
			"password too short",
			slog.Int("length", passLength))
		httputil.RespondWithError(w, r, http.StatusBadRequest, "password too short, length less than 8")
		return
	}

	if params.HashedPassword, err = auth.HashPassword(params.HashedPassword); err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"an error occured when hashind the password",
			slog.String("err", err.Error()),
			slog.String("user email", params.Email))
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "password hashing error")
		return
	}

	usr, err := app.DB.CreateUser(r.Context(), params)
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

		httputil.RespondWithError(w, r, status, msg)
		return
	}

	httputil.RespondWithJSON(w, r, http.StatusCreated, usr)
}
