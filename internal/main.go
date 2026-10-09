package main

import (
	"fmt"
	"kogatari/internal/handlers"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("an error has occured when accessing env variables", err)
	}

	srvPort := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/test", handlers.RootHandler)

	srv := &http.Server{
		Addr:    ":" + srvPort,
		Handler: mux,
	}

	fmt.Println("Server running on port:", srvPort)
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("an error has occured when running server", err)
	}
}
