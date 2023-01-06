package main

import (
	"bytes"
	"encoding/json"
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
			msg, err := deserialize(sms.Body)
			if err != nil {
				panic(err)
			}
			log.Println("Received SMS:", msg.Message)
		}
	}()

	<-forever
}

type Message map[string]interface{}

func deserialize(b []byte) (helpers.GroupMessage, error) {
	// speel met if en reflect operator om te bepalen of het een group of location message is
	var msg helpers.GroupMessage
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
