package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Global variables for RabbitMQ connection and channel
var (
	channel *amqp.Channel
	conn    *amqp.Connection
)

// Exchange and Queue constants
const (
	// Work Exchange and Queue for initial message processing
	WORK_EXCHANGE_NAME = "GPcoder.WorkExchange"
	WORK_QUEUE_NAME    = "WorkQueue"

	// Retry Exchange and Queue for dead-letter processing
	RETRY_EXCHANGE_NAME = "GPcoder.RetryExchange"
	RETRY_QUEUE_NAME    = "RetryQueue"

	// Retry delay for messages in the dead-letter queue (in milliseconds)
	RETRY_DELAY = 300
)

// A simple global counter to track retry attempts.
// NOTE: In a real-world application, this would be a race condition with multiple
// consumers. A robust solution would use a message header to track the retry count.
var retryCount int

func main() {
	// Establish a connection to the RabbitMQ server.
	// The host is set to "localhost", which can be changed.
	factory := amqp.DialConfig{
		Dial:      amqp.Dial,
		TLSConfig: nil,
	}
	var err error
	conn, err = factory.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", "localhost"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel for communication.
	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	// Declare the Work Exchange and Queue.
	// The Work Queue is declared with no special arguments.
	err = channel.ExchangeDeclare(WORK_EXCHANGE_NAME, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare Work Exchange: %v", err)
	}
	_, err = channel.QueueDeclare(WORK_QUEUE_NAME, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare Work Queue: %v", err)
	}
	err = channel.QueueBind(WORK_QUEUE_NAME, "", WORK_EXCHANGE_NAME, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind Work Queue: %v", err)
	}

	// Declare the Retry Exchange and Queue (the Dead Letter Queue).
	// This queue is configured with a dead-letter exchange and a message TTL.
	// Messages that expire in this queue are routed to the dead-letter exchange,
	// which is the Work Exchange.
	args := amqp.Table{
		"x-dead-letter-exchange": WORK_EXCHANGE_NAME,
		"x-message-ttl":          RETRY_DELAY,
	}
	err = channel.ExchangeDeclare(RETRY_EXCHANGE_NAME, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare Retry Exchange: %v", err)
	}
	_, err = channel.QueueDeclare(RETRY_QUEUE_NAME, true, false, false, false, args)
	if err != nil {
		log.Fatalf("Failed to declare Retry Queue: %v", err)
	}
	err = channel.QueueBind(RETRY_QUEUE_NAME, "", RETRY_EXCHANGE_NAME, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind Retry Queue: %v", err)
	}

	// Publish an initial message to the Work Exchange.
	message := "GPCoder Message"
	log.Printf("[%s] [Work] [Send]: %s", time.Now().Format("15:04:05"), message)
	err = channel.Publish(WORK_EXCHANGE_NAME, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	// Start the consumer. This function will block and listen for messages.
	// The `true` argument here enables auto-acknowledgement.
	consumer(WORK_QUEUE_NAME)
}

// consumer listens for and processes messages from the specified queue.
func consumer(queueName string) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Wait for messages indefinitely.
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			content := string(d.Body)
			log.Printf("[%s] [Received] [%s]: %s\n", time.Now().Format("15:04:05"), queueName, content)

			// Logic to decide whether to retry the message.
			if retryCount < 5 {
				publishToRetryExchange(content)
				retryCount++
			} else {
				// After 5 retries, reset the counter.
				// NOTE: As a result, the message will be retried indefinitely.
				// A more robust system would stop retrying after a certain count.
				retryCount = 0
			}
		}
	}()

	log.Printf("Waiting for messages from queue: %s. To exit press CTRL+C", queueName)
	<-forever
}

// publishToRetryExchange sends the message to the retry exchange.
func publishToRetryExchange(message string) {
	log.Printf("[%s] [Retry%d] [Re-Publish]: %s", time.Now().Format("15:04:05"), retryCount+1, message)
	err := channel.Publish(
		RETRY_EXCHANGE_NAME, // exchange
		"",                  // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish to retry exchange: %v", err)
	}
}
