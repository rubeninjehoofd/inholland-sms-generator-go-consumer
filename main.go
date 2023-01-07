package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sms-consumer/app/helpers"

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

			if gMsg.ClassId != "" {
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

// deserialize the byte array to json message object
func deserializeToJson(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

// Determines the message type
func determineMessageType(msg Message) (helpers.GroupMessageJSON, helpers.LocationMessageJSON) {
	var gMsg helpers.GroupMessageJSON
	var lMsg helpers.LocationMessageJSON
	if val, ok := msg["ClassId"]; ok {
		fmt.Println(reflect.TypeOf(val))
		gMsg = deserializeToGroupMessage(msg)
	} else {
		lMsg = deserializeToLocationMessage(msg)
	}
	return gMsg, lMsg
}

// Makes a group message object with only string field from the json message
// (uuid, time, etc. types are not in the json object)
func deserializeToGroupMessage(msg Message) helpers.GroupMessageJSON {
	gMsg := helpers.GroupMessageJSON{
		MessageId:       msg["MessageId"].(string),
		ClassId:         msg["ClassId"].(string),
		ScheduledAt:     msg["ScheduledAt"].(string),
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}
	return gMsg
}

// Makes a location message object with only string field from the json message
// (uuid, time, etc. types are not in the json object)
func deserializeToLocationMessage(msg Message) helpers.LocationMessageJSON {
	lMsg := helpers.LocationMessageJSON{
		MessageId:       msg["MessageId"].(string),
		LocationId:      msg["LocationId"].(string),
		ScheduledAt:     msg["ScheduledAt"].(string),
		Message:         msg["Message"].(string),
		FromPhoneNumber: msg["FromPhoneNumber"].(string),
		ToPhoneNumber:   msg["ToPhoneNumber"].(string),
	}
	return lMsg
}
