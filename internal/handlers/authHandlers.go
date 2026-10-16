package handlers

import (
	"kogatari/internal/auth"
	"kogatari/internal/database"
	"kogatari/internal/httputil"
	"log/slog"
	"net/http"
)

func (app application) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	}

	refreshTokenDB, err := app.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		app.Slogger.WarnContext(r.Context(),
			"error when accessing the refresh tokens",
			slog.String("err", err.Error()))
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	} else if isValid, err := auth.ValidateRefreshToken(refreshTokenDB); !isValid {
		app.Slogger.WarnContext(r.Context(),
			"error when accessing the refresh tokens",
			slog.String("err", err.Error()))
		httputil.RespondWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	freshToken := auth.MakeRefreshToken()
	err = app.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:   refreshTokenDB.Token,
		Token_2: freshToken,
	})
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	httputil.RespondWithJSON(w, r, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: freshToken,
	})
}
