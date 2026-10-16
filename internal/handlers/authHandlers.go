package handlers

import (
	"kogatari/internal/auth"
	"kogatari/internal/httputil"
	"net/http"
	"time"
)

func (app application) RefreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	}

	refreshTokenDB, err := app.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	} else if isValid, err := auth.ValidateRefreshToken(refreshTokenDB); !isValid {
		httputil.RespondWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	freshToken, err := auth.GenerateJWT(refreshTokenDB.UserID, app.JWTKey, time.Hour)
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

func (app application) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	}

	refreshTokenDB, err := app.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusBadRequest, "invalid header")
		return
	} else if isValid, err := auth.ValidateRefreshToken(refreshTokenDB); !isValid {
		httputil.RespondWithError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	err = app.DB.RevokeRefreshToken(r.Context(), refreshTokenDB.Token)
	if err != nil {
		httputil.RespondWithError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
