package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	Message   []byte    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
