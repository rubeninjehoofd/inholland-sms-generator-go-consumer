package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var (
	accountSid string
	authToken  string
	fromPhone  string
	toPhone    string
	client     *twilio.RestClient
)

func SendMessage(msg string) {

	params := openapi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(fromPhone)
	params.SetBody(msg)

	response, err := client.Api.CreateMessage(&params)
	if err != nil {
		fmt.Printf("error creating and sending message: %s\n", err.Error())
		return
	}
	fmt.Printf("Message SID: %s\n", *response.Sid)
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env: %s\n", err.Error())
		os.Exit(1)
	}

	accountSid = "AC7c1bac13c8d94ff8b152a06f74e47d47"
	authToken = "e3c1d0d2aa98708bf66689a17251c88c"
	fromPhone = "+12183199239"
	toPhone = "+31611523882"

	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
}

func sender() {

	msg := fmt.Sprintf(os.Getenv("MSG"), "James")
	SendMessage(msg)

}
