package entities

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
