package main

import (
	"fmt"
	"sms-consumer/src/helpers"
	"time"

	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	accountSid string
	authToken  string
	client     *twilio.RestClient
)

func SendMessage(msg helpers.GroupMessage) {
	params := openapi.CreateMessageParams{}
	params.SetMessagingServiceSid("{MessagingServiceID}")

	// check if scheduled
	if msg.ScheduledAt.Unix() > (time.Now().Local().Unix() + 900) {
		fmt.Println("schedule sms")
		params.SetSendAt(msg.ScheduledAt)
		params.SetScheduleType("fixed")
	}

	// sms params
	params.SetTo(msg.ToPhoneNumber)
	params.SetFrom(msg.FromPhoneNumber)
	params.SetBody(msg.Message)

	// send sms
	resp, err := client.Api.CreateMessage(&params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}
}

func init() {
	accountSid = "{AccountSid}"
	authToken = "{authToken}"

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

func sender() {

	var GroupMessage helpers.GroupMessage
	GroupMessage.MessageId = uuid.UUID{}
	GroupMessage.ClassId = uuid.UUID{}
	//GroupMessage.ScheduledAt = time.Now()
	GroupMessage.ScheduledAt = time.Date(2023, 1, 6, 15, 18, 0, 0, time.Local)
	GroupMessage.Message = "test"
	GroupMessage.FromPhoneNumber = "{FromPhoneNumber}"
	GroupMessage.ToPhoneNumber = "{ToPhoneNumber}"
	SendMessage(GroupMessage)

	// var LocationMessage helpers.LocationMessage
	// LocationMessage.MessageId = uuid.UUID{}
	// LocationMessage.LocationId = uuid.UUID{}
	// LocationMessage.ScheduledAt = time.Now()
	// LocationMessage.ScheduledAt = time.Date(2023, 1, 6, 15, 18, 0, 0, time.Local)
	// LocationMessage.Message = "test"
	// LocationMessage.FromPhoneNumber = "{FromPhoneNumber}"
	// LocationMessage.ToPhoneNumber = "{ToPhoneNumber}"
	//SendMessage(LocationMessage)

}
