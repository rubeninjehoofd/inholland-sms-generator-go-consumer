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
