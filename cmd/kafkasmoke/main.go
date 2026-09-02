package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/budaevdev/roomly/internal/kafka"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	brokers := os.Getenv("KAFKA_BROKERS")

	consumer := kafka.NewConsumer([]string{brokers}, "test-topic", "smoketest-group")
	defer consumer.Close()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		value, err := consumer.ReadMessage(ctx)
		if err != nil {
			fmt.Println("consumer error:", err)
			return
		}
		fmt.Println("consumed:", string(value))
	}()

	time.Sleep(2 * time.Second)

	producer := kafka.NewProducer([]string{brokers})
	defer producer.Close()

	err := producer.Publish("test-topic", []byte("hello"), []byte("world"))
	if err != nil {
		fmt.Println("publish error:", err)
		return
	}
	fmt.Println("published")

	time.Sleep(11 * time.Second)

}
