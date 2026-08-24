package main

import (
	"log"
	"net/http"
	"os"

	"github.com/budaevdev/roomly/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	conn, err := storage.Connect(dbURL)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, mux))
}
