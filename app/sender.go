package app

import (
	"fmt"
	"os"
	"sms-consumer/app/helpers"
	"time"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	accountSid        string
	authToken         string
	client            *twilio.RestClient
	twilioPhoneNumber string
	messageServiceId  string
)

func SendMessage(msg helpers.BaseMessage) {
	params := openapi.CreateMessageParams{}
	params.SetMessagingServiceSid(messageServiceId)

	// check if scheduled
	if msg.ScheduledAt.Unix() > (time.Now().Local().Unix() + 900) {
		fmt.Println("schedule sms")
		params.SetSendAt(msg.ScheduledAt)
		params.SetScheduleType("fixed")
	}

	// sms params
	params.SetTo(msg.ToPhoneNumber)
	params.SetFrom(twilioPhoneNumber)
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

func Init() {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	twilioPhoneNumber = os.Getenv("TWILIO_NUMBER")
	messageServiceId = os.Getenv("MESSAGE_SERVICE_ID")

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}
