package helpers

import (
	"github.com/google/uuid"
)

type GroupMessage struct {
	BaseMessage
	ClassId uuid.UUID
}
