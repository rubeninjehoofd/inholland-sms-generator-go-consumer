package helpers

import (
	"time"
)

type BaseMessage struct {
	ScheduledAt     time.Time
	Message         string
	FromPhoneNumber string
	ToPhoneNumber   string
}
