package entities

import "time"

type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	WIPLimit  *int32    `json:"wip_limit,omitempty"`
	Position  float64   `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
