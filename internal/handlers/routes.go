package handlers

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/register", app.CreateUserHander)
	mux.HandleFunc("/healthz", HealthHandler)
}
