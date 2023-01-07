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
