package entities

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
