package main

import (
	"context"
	"kogatari/internal/handlers"
	"log/slog"
)

func main() {
	srv, dbConn := handlers.CreateServer()

	defer dbConn.Close(context.Background())

	slog.Info("started server")
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("an error has occured when running server",
			slog.String("err", err.Error()))
	}
}
