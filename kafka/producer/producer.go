package main

import (
	"log"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // Ensures you can track delivery
	config.Producer.Retry.Max = 5           // Retry for transient failures

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder("Hello Kafka!"),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
