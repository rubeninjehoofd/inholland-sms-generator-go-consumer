package helpers

import (
	"time"

	"github.com/google/uuid"
)

type GroupMessage struct {
	MessageId       uuid.UUID
	ClassId         uuid.UUID
	ScheduledAt     time.Time
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}

// There is a json version, because once the byte array is deserialized, all
// the fields in that object are of type string
type GroupMessageJSON struct {
	MessageId       string
	ClassId         string
	ScheduledAt     string
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}
