package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("SMS Consumer - Connecting to the SMS channel")
	log.Println("SMS Consumer - Connecting to the SMS channel")

	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://guest:guest@localhost:5672/"

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
