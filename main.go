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
			deserializedMsg, err := deserializeToJson(sms.Body)
			if err != nil {
				panic(err)
			}
			gMsg, lMsg := determineMessageType(deserializedMsg)

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

// Deserialize the byte array to json message object
func deserializeToJson(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

// Determines the message type
func determineMessageType(msg Message) (helpers.GroupMessage, helpers.LocationMessage) {
	var gMsg helpers.GroupMessage
	var lMsg helpers.LocationMessage
	if _, ok := msg["ClassId"]; ok {
		gMsg = deserializeToGroupMessage(msg)
	} else {
		lMsg = deserializeToLocationMessage(msg)
	}
	return gMsg, lMsg
}

// Makes a group message object with only string field from the json message
// (uuid, time, etc. types are not in the json object)
func deserializeToGroupMessage(msg Message) helpers.GroupMessage {
	// Parse can't happen inline, because the time.Parse function can return 2 values
	scheduledAt, err := time.Parse(time.RFC3339, msg["ScheduledAt"].(string))
	if err != nil {
		panic(err)
	}

	// Make group message
	gMsg := helpers.GroupMessage{
		MessageId:       uuid.MustParse(msg["MessageId"].(string)),
		ClassId:         uuid.MustParse(msg["ClassId"].(string)),
		ScheduledAt:     scheduledAt,
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}
	return gMsg
}

// Makes a location message object with only string field from the json message
// (uuid, time, etc. types are not in the json object)
func deserializeToLocationMessage(msg Message) helpers.LocationMessage {
	// Parse can't happen inline, because the time.Parse function can return 2 values
	scheduledAt, err := time.Parse(time.RFC3339, msg["ScheduledAt"].(string))
	if err != nil {
		panic(err)
	}

	// Make location message
	lMsg := helpers.LocationMessage{
		MessageId:       uuid.MustParse(msg["MessageId"].(string)),
		LocationId:      uuid.MustParse(msg["LocationId"].(string)),
		ScheduledAt:     scheduledAt,
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}
	return lMsg
}
