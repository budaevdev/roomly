package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/budaevdev/roomly/internal/booking"
	"github.com/budaevdev/roomly/internal/cache"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb, err := cache.Connect(redisAddr)

	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

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

	mux.HandleFunc("GET /listings/search", func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")

		if start == "" || end == "" {
			http.Error(w, "start and end query params are required", http.StatusBadRequest)
			return
		}

		key := fmt.Sprintf("avail:%s:%s", start, end)

		val, err := rdb.Get(r.Context(), key).Result()

		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(val))
			return
		}

		listings, err := storage.GetAvailableListings(conn, start, end)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(listings)
		if err == nil {
			rdb.Set(r.Context(), key, data, 30*time.Second)
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

	mux.HandleFunc("PATCH /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = storage.CancelBooking(conn, id)

		if err != nil {
			if errors.Is(err, storage.ErrBookingNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	port := ":8080"
	log.Fatal(http.ListenAndServe(port, mux))
}
