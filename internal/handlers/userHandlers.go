package handlers

import (
	"encoding/json"
	"kogatari/internal/auth"
	"kogatari/internal/database"
	"kogatari/internal/httputil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"uuid"
)

func (app application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params database.CreateUserParams
	var err error

	if err = decoder.Decode(&params); err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	} else if len(params.Email) == 0 || len(params.HashedPassword) < 8 {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	}
	params.Email = strings.ToLower(params.Email)

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
		app.Slogger.ErrorContext(r.Context(),
			"error with creating user has occured",
			slog.String("err", err.Error()),
			slog.String("email", params.Email))
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "internal error")
		return
	}

	httputil.RespondWithJSON(w, r, http.StatusCreated, usr)
}

func (app application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Email            string `json:"email"`
		HashedPassword   string `json:"hashed_password"`
		ExpiresInSeconds string `json:"expires_in_seconds"`
	}
	var err error
	if err = decoder.Decode(&params); err != nil ||
		len(params.Email) == 0 ||
		len(params.HashedPassword) == 0 {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "bad body request")
		return
	}
	params.Email = strings.ToLower(params.Email)

	usr, err := app.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "couldnt retrieve user data")
		return
	}

	if passwordMatch, err := auth.CheckPasswordHash(params.HashedPassword, usr.HashedPassword); err != nil {
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "password hash check error")
		return
	} else if !passwordMatch {
		app.Slogger.WarnContext(r.Context(),
			"failed log in attempt",
			slog.String("user email", params.Email))
		httputil.RespondWithError(w, r, http.StatusUnauthorized, "invalid email or password")
		return
	}

	expiryTime := time.Hour * 1
	if params.ExpiresInSeconds != "" {
		multiplier, err := strconv.Atoi(params.ExpiresInSeconds)
		if err != nil {
			httputil.RespondWithError(w, r, http.StatusInternalServerError, "internal error")
		}
		expiryTime = time.Second * time.Duration(multiplier)
	}
	token, err := auth.GenerateJWT(usr.ID, app.JWTKey, expiryTime)
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "internal error")
	}

	httputil.RespondWithJSON(w, r, http.StatusOK, struct {
		Email string    `json:"email"`
		UUID  uuid.UUID `json:"uuid"`
		Token string    `json:"token"`
	}{
		Email: usr.Email,
		UUID:  uuid.UUID(usr.ID),
		Token: token,
	})
}
