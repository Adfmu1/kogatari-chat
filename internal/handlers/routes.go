package handlers

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	// POST
	mux.HandleFunc("POST /api/v1/register", app.CreateUserHandler)

	// GET
	mux.HandleFunc("GET /api/v1/login", app.LoginUserHandler)
	mux.HandleFunc("GET /api/v1/refresh", app.RefreshJWTHandler)
	mux.HandleFunc("GET /api/v1/revoke", app.RevokeRefreshTokenHandler)
	mux.HandleFunc("GET /healthz", HealthHandler)
}
