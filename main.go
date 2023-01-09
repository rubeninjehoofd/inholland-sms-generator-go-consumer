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

	// Define RabbitMQ server URL.
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
			deserializedMsg := deserializeToJson(sms.Body)

			gMsg, lMsg := createMessage(deserializedMsg)

			if gMsg.ClassId != uuid.Nil {
				// sender.SendGroupMessage(gMsg)
				log.Println("Received class message:", gMsg.Message)

			} else {
				// sender.SendLocationMessage(lMsg)
				log.Println("Received location message:", lMsg.Message)
			}
		}
	}()

	<-forever
}

type Message map[string]interface{}

// Deserializes the byte array to a json message object
func deserializeToJson(b []byte) Message {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	if err != nil {
		panic(err)
	}
	return msg
}

// Creates a Location message or Group message, based on
// the contents of the msg variable
func createMessage(msg Message) (helpers.GroupMessage, helpers.LocationMessage) {
	var gMsg helpers.GroupMessage
	var lMsg helpers.LocationMessage

	// Declare the base message
	baseMsg := helpers.BaseMessage{
		MessageId:       uuid.MustParse(msg["MessageId"].(string)),
		ScheduledAt:     parseTime(msg),
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}

	// Check what the message type is
	if _, ok := msg["ClassId"]; ok {
		// Create Group message
		gMsg = helpers.GroupMessage{
			BaseMessage: baseMsg,
			ClassId:     uuid.MustParse(msg["ClassId"].(string)),
		}
	} else {
		// Create location message
		lMsg = helpers.LocationMessage{
			BaseMessage: baseMsg,
			LocationId:  uuid.MustParse(msg["LocationId"].(string)),
		}
	}
	return gMsg, lMsg
}

// Parses the string time to time.Time
func parseTime(msg Message) time.Time {
	scheduledAt, err := time.Parse(time.RFC3339, msg["ScheduledAt"].(string))
	if err != nil {
		panic(err)
	}
	return scheduledAt
}
