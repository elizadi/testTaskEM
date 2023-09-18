package kafka

import (
	"context"
	"effective_mobile/internal/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type repository struct {
	writer *kafka.Writer
	reader *kafka.Reader
	// writerTopic string
}

// todo
func New() (domain.KafkaRepository, error) {
	port1 := os.Getenv("kafka_port1")
	if port1 == "" {
		log.Fatalf("set kafka_port1 in env file")
	}
	port2 := os.Getenv("kafka_port2")
	if port2 == "" {
		log.Fatalf("set kafka_port2 in env file")
	}
	port3 := os.Getenv("kafka_port3")
	if port3 == "" {
		log.Fatalf("set kafka_port3 in env file")
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(port1, port2),
		Topic:    "FIO_FAILED",
		Logger:   logrus.New(),
		Balancer: &kafka.LeastBytes{},
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          []string{port1, port2},
		GroupID:          "this_service",
		Topic:            "FIO",
		Logger:           logrus.New(),
		Dialer:           kafka.DefaultDialer,
		ReadBatchTimeout: time.Second,
		MinBytes:         10e3, // same value of Shopify/sarama
		MaxBytes:         10e6,
	})
	c, err := kafka.Dial("tcp", port3)
	if err != nil {
		c.Close()
	}
	kt := kafka.TopicConfig{Topic: "FIO_FAILED", NumPartitions: 2, ReplicationFactor: 2}
	err = c.CreateTopics(kt)
	if err != nil {
		log.Println(err)
	}
	kt = kafka.TopicConfig{Topic: "FIO", NumPartitions: 2, ReplicationFactor: 2}
	err = c.CreateTopics(kt)
	if err != nil {
		log.Println(err)
	}

	return &repository{
		writer: w,
		reader: r,
		// writerTopic: "FIO_FAILED",
	}, nil
}

func (r *repository) Produce(message string) error {
	return r.writer.WriteMessages(context.Background(), kafka.Message{
		// Topic: r.writerTopic,
		Value: []byte(message),
	})
}

func (r *repository) Consume() (domain.FIOUser, error) {
	ctx, cance := context.WithTimeout(context.Background(), time.Second*5)
	defer cance()
	message, err := r.reader.ReadMessage(ctx)
	if err != nil && err != io.ErrUnexpectedEOF && err != context.Canceled && err != context.DeadlineExceeded {
		return domain.FIOUser{}, err
	}
	if (err == io.ErrUnexpectedEOF || err == context.DeadlineExceeded) && len(message.Value) == 0 {
		return domain.FIOUser{}, domain.ErrEmptyMessage
	}
	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
		message.Topic,
		message.Partition,
		message.Offset,
		string(message.Key),
		string(message.Value))
	var fio domain.FIOUser
	err = json.Unmarshal(message.Value, &fio)
	if err != nil {
		return domain.FIOUser{}, err
	}
	return fio, nil
}
