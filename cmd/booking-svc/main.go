package main

import (
	"log"
	"os"

	"github.com/budaevdev/roomly/internal/api"
	"github.com/budaevdev/roomly/internal/cache"
	"github.com/budaevdev/roomly/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
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

	srv := api.NewServer(conn, rdb)

	log.Fatal(srv.Routes().Run(":8080"))
}
