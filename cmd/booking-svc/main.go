package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/budaevdev/roomly/internal/booking"
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

	mux.HandleFunc("POST /listings", func(w http.ResponseWriter, r *http.Request) {
		var l booking.Listing

		err := json.NewDecoder(r.Body).Decode(&l)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = storage.CreateListing(conn, &l)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(l)
	})

	mux.HandleFunc("GET /listings", func(w http.ResponseWriter, r *http.Request) {
		listings, err := storage.GetListings(conn)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(listings)
	})

	mux.HandleFunc("POST /bookings", func(w http.ResponseWriter, r *http.Request) {
		var b booking.Booking

		err := json.NewDecoder(r.Body).Decode(&b)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = storage.CreateBooking(conn, &b)

		if err != nil {
			if errors.Is(err, storage.ErrOverlappingBooking) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)
	})

	port := ":8080"
	log.Fatal(http.ListenAndServe(port, mux))
}
