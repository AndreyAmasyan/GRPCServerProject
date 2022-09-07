package kafkaBroker

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"test/internal/data"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic     string = "data_topic"
	partition int    = 0
)

// KafkaProducer ...
type KafkaProducer struct {
	conn *kafka.Conn
}

func New(ctx context.Context) (*KafkaProducer, error) {

	p, err := kafka.DefaultDialer.LookupPartition(ctx, "tcp", os.Getenv("KAFKA_HOST_PORT"), topic, partition)
	if err != nil {
		fmt.Println("KafkaProducer New kafka.DefaultDialer.LookupPartition", err)
		return nil, err
	}
	p.Leader.Host = strings.Split(os.Getenv("KAFKA_HOST_PORT"), ":")[0]
	conn, err := kafka.DefaultDialer.DialPartition(ctx, "tcp", os.Getenv("KAFKA_HOST_PORT"), p)
	if err != nil {
		fmt.Println("KafkaProducer New kafka.DefaultDialer.DialPartition", err)
		return nil, err
	}
	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return &KafkaProducer{
		conn: conn,
	}, nil
}

func (s *KafkaProducer) Produce(user data.User) error {
	n, err := s.conn.WriteMessages(kafka.Message{
		Value: []byte("\"" + strconv.Itoa(user.ID) + "\",\"" + user.UserName + "\",\"" + time.Now().String() + "\"\r\n"),
	})
	if err != nil {
		fmt.Println("KafkaProducer Produce s.conn.WriteMessages", err)
		return err
	}
	fmt.Println(n, " bytes written (Kafka Produce)")

	return nil
}

func (s *KafkaProducer) Close() error {
	if err := s.conn.Close(); err != nil {
		fmt.Println("KafkaProducer Close s.conn.Close", err)
		return err
	}

	return nil
}
