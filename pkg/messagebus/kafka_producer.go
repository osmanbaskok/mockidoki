package messagebus

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"mockidoki/config"
	"strconv"
	"time"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config config.KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Broker),
		RequiredAcks: kafka.RequireAll,
	}

	return &KafkaProducer{writer: writer}
}

func (producer *KafkaProducer) Produce(message []byte, channel string) error {
	producer.writer.Topic = channel

	nSec := time.Now().UnixNano()

	err := producer.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.FormatInt(nSec, 10)),
		Value: message,
	})

	if err != nil {
		return err
	}
	return nil
}

func (producer *KafkaProducer) ProduceList(messageList []interface{}, channel string) error {
	producer.writer.Topic = channel

	nSec := time.Now().UnixNano()

	messages := make([]kafka.Message, 0)
	for _, message := range messageList {
		eventMessage, _ := json.Marshal(message)

		kafkaMessage := kafka.Message{
			Key:   []byte(strconv.FormatInt(nSec, 10)),
			Value: eventMessage,
		}
		messages = append(messages, kafkaMessage)
	}

	err := producer.writer.WriteMessages(context.Background(), messages...)

	if err != nil {
		return err
	}
	return nil
}
