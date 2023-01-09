package helpers

import (
	"github.com/google/uuid"
)

type LocationMessage struct {
	BaseMessage
	LocationId uuid.UUID
}
