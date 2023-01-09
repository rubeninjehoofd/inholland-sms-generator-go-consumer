package helpers

import (
	"time"

	"github.com/google/uuid"
)

type BaseMessage struct {
	MessageId       uuid.UUID
	ScheduledAt     time.Time
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}
