package handlers

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	// POST
	mux.HandleFunc("POST /api/v1/register", app.CreateUserHander)

	// GET
	mux.HandleFunc("GET /api/v1/login", app.LoginUserHandler)
	mux.HandleFunc("GET /healthz", HealthHandler)
}
