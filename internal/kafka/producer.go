package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	return &Producer{writer: w}
}

func (p *Producer) Publish(topic string, key, value []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
