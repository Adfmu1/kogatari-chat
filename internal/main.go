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
	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("an error has occured when running server",
			slog.String("err", err.Error()))
	}
}
