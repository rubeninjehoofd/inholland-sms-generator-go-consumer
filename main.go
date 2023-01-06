package main

import (
	"fmt"
	"log"
	"sms-consumer/app/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("SMS Consumer - Connecting to the SMS channel")
	log.Println("SMS Consumer - Connecting to the SMS channel")

	classMsg := helpers.GroupMessage{
		MessageId:   uuid.New(),
		ClassId:     uuid.New(),
		ScheduledAt: time.Now(),
		Message:     "hello world",
		PhoneNumber: "test",
	}

	fmt.Println(classMsg.Message)

	// Define RabbitMQ server URL.
	// amqpServerURL := os.Getenv("AMQP_SERVER_URL_TEST")
	amqpServerURL := "amqp://guest:guest@rabbitmq:5672/"

	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	smsChannel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer smsChannel.Close()

	messages, err := smsChannel.Consume(
		"SMSQueue", // queue name
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to SMS channel")
	log.Println("Waiting for messages")

	forever := make(chan bool)

	go func() {
		for sms := range messages {
			log.Println("Received SMS:", sms.Body)
		}
	}()

	<-forever
}
