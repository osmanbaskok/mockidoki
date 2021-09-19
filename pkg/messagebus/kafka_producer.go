package messagebus

import (
	"context"
	"github.com/segmentio/kafka-go"
	"mockidoki/config"
	"strconv"
	"time"
)

type KafkaProducer struct {
	broker string
}

func NewKafkaProducer(config config.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{broker: config.Broker}
}

func (producer *KafkaProducer) Produce(message string, channel string) error {
	w := &kafka.Writer{
		Addr:         kafka.TCP(producer.broker),
		Topic:        channel,
		RequiredAcks: kafka.RequireAll,
	}

	nSec := time.Now().UnixNano()

	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.FormatInt(nSec, 10)),
		Value: []byte(message),
	})

	if err != nil {
		return err
	}
	return nil
}
