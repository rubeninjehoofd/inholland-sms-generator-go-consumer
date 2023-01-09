package app

import (
	"fmt"
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

func SendBaseMessage(msg helpers.BaseMessage) {
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

			// TODO: this wil send confirmation for every sms (teacher gets spammed)
			SendConfirmationToTeacher(msg)
		} else {
			fmt.Println(resp.Sid)
		}
	}
}

func SendConfirmationToTeacher(msg helpers.BaseMessage) {
	params := openapi.CreateMessageParams{}
	params.SetMessagingServiceSid(messageServiceId)

	// check if scheduled
	if msg.ScheduledAt.Unix() > (time.Now().Local().Unix() + 900) {
		fmt.Println("schedule sms")
		params.SetSendAt(msg.ScheduledAt)
		params.SetScheduleType("fixed")
	}

	// sms params
	params.SetTo(msg.FromPhoneNumber)
	params.SetFrom(twilioPhoneNumber)
	params.SetBody("SMS confimation - Message has been sent")

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
	accountSid = "AC7c1bac13c8d94ff8b152a06f74e47d47"
	authToken = "e3c1d0d2aa98708bf66689a17251c88c"
	twilioPhoneNumber = "+12183199239"
	messageServiceId = "MGc04be72a20e29d49ecd83f987720b9c4"

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}
