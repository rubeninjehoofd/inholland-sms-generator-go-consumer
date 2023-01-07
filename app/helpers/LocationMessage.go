package helpers

import (
	"time"

	"github.com/google/uuid"
)

type LocationMessage struct {
	MessageId       uuid.UUID
	LocationId      uuid.UUID
	ScheduledAt     time.Time
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}

// There is a json version, because once the byte array is deserialized, all
// the fields in that object are of type string
type LocationMessageJSON struct {
	MessageId       string
	LocationId      string
	ScheduledAt     string
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}
