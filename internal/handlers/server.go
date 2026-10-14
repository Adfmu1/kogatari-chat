package handlers

import (
	"context"
	"kogatari/internal/database"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type application struct {
	DB      *database.Queries
	Slogger *slog.Logger
}

var app application

func CreateServer() (*http.Server, *pgx.Conn) {
	app.Slogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		app.Slogger.Error("error loading env variables",
			slog.String("err", err.Error()))
		return nil, nil
	}

	srvPort := os.Getenv("PORT")
	dbURL := os.Getenv("DB_STRING")

	dbConn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		app.Slogger.Error("error with connecting to database",
			slog.String("error", err.Error()))
		return nil, nil
	}

	app.DB = database.New(dbConn)
	app.Slogger.Info("connected to db")

	mux := http.NewServeMux()

	registerRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + srvPort,
		Handler: mux,
	}

	return srv, dbConn
}
