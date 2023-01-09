package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	sender "sms-consumer/app"
	"sms-consumer/app/helpers"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	log.Println("SMS Consumer - Connecting to the SMS channel")

	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

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
	sender.Init()

	go func() {
		for sms := range messages {
			// Prepare the message
			msg := prepareMessage(sms.Body)

			// Send the message
			sender.SendBaseMessage(msg)
		}
	}()

	<-forever
}

type Message map[string]interface{}

// Deserializes the byte array to a json message object
func prepareMessage(bytes []byte) helpers.BaseMessage {
	return createMessage(deserialize(bytes))
}

// Deserializes the byte array to a json message object
func deserialize(b []byte) Message {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	if err != nil {
		panic(err)
	}
	return msg
}

// Creates a message
func createMessage(msg Message) helpers.BaseMessage {
	return helpers.BaseMessage{
		ScheduledAt:     parseTime(msg),
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}
}

// Parses the string time to time.Time
func parseTime(msg Message) time.Time {
	scheduledAt, err := time.Parse(time.RFC3339, msg["ScheduledAt"].(string))
	if err != nil {
		log.Println("Not scheduled for a specific time. Message will be send now")
		return time.Now()
	}
	return scheduledAt
}
