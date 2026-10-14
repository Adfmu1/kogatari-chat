package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"kogatari/internal/auth"
	"kogatari/internal/database"
	"kogatari/internal/httputil"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func (app application) CreateUserHander(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params database.CreateUserParams
	var err error

	if err = decoder.Decode(&params); err != nil {
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

func (app application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var err error
	if err = decoder.Decode(&params); err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"error with request body has occured",
			slog.String("err", err.Error()))

		httputil.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	}

	params.Email = strings.ToLower(params.Email)

	hash, err := app.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		msg := "couldnt retrieve user data"
		status := http.StatusInternalServerError
		if errors.Is(err, sql.ErrNoRows) {
			msg = "no user with this email"
			status = http.StatusUnauthorized
		}
		app.Slogger.ErrorContext(r.Context(),
			"error with logging in user has occured",
			slog.String("err", err.Error()),
			slog.String("email", params.Email))
		httputil.RespondWithError(w, r, status, msg)
		return
	}

	if passwordMatch, err := auth.CheckPasswordHash(params.Password, hash); err != nil {
		app.Slogger.ErrorContext(r.Context(),
			"an error occured when checking hashed password",
			slog.String("err", err.Error()),
			slog.String("user email", params.Email))
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "password hash check error")
		return
	} else if !passwordMatch {
		app.Slogger.ErrorContext(r.Context(),
			"failed log in attempt",
			slog.String("user email", params.Email))
		httputil.RespondWithError(w, r, http.StatusUnauthorized, "passwords don't match")
		return
	}

	w.WriteHeader(http.StatusOK)
}
